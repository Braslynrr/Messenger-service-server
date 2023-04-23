package messengerserviceapi

import (
	"MessengerService/group"
	messengermanager "MessengerService/mesermanager"
	msmessage "MessengerService/message"
	"MessengerService/user"
	"MessengerService/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ms *MessengerService) newSocketIo() http.Handler {
	ms.socketIO = socket.NewServer(nil, nil)

	ms.socketIO.On("connection", func(clients ...any) {
		var err error

		client := clients[0].(*socket.Socket)

		ms.Logger.Printf("New Connection: %v", client.Id())

		if client.Handshake().Query != nil {
			token, ok := client.Handshake().Query.Get("token")
			if ok {

				WSKey, _ := utils.GenerateRandomAESKey(16)
				ms.Logger.Printf("Sending key: %v", client.Id())
				client.Emit("WSKey", WSKey)
				client.SetData(WSKey)

				if ms.connectToMessengerService(client, token) {
					ms.sockets[client.Id()] = client

					// options
					client.On("sendMessage", func(message ...any) {

						ms.handleError(client, "", hasUserMiddleware(client, ms.sendMessage, message...))
					})

					client.On("GetCurrentGroups", func(...any) {
						ms.handleError(client, "", hasUserMiddlewareNoParam(client, ms.getGroups))
					})

					client.On("GroupHistory", func(history ...any) {
						ms.handleError(client, "", hasUserMiddleware(client, ms.getGroupHistory, history...))
					})

					client.On("SendSeen", func(SeenMesage ...any) {
						ms.handleError(client, "", hasUserMiddleware(client, ms.seenMessage, SeenMesage...))
					})

					client.On("CreateGroup", func(group ...any) {
						ms.handleError(client, "", hasUserMiddleware(client, ms.createGroup, group...))
					})
					client.On("disconnect", func(reason ...any) {
						ms.Logger.Printf("%s is disconnecting because %v ", client.Id(), reason[0])
						client.Disconnect(true)
						delete(ms.sockets, client.Id())
					})
				} else {
					err = errors.New("could not connect, invalid token")
				}

			} else {
				err = errors.New("token was not sent")
			}

		} else {
			err = errors.New("it's necessary send a query")
		}

		if err != nil {
			ms.handleError(client, "", err)
			client.Disconnect(true)
		}

	})

	return ms.socketIO.ServeHandler(nil)
}

// handleError handles erros using a type and error
func (ms *MessengerService) handleError(conn *socket.Socket, errortype string, err error) {
	if err != nil {
		ms.ErrorChan <- SocketError{err: err, socket: conn, errorType: fmt.Sprintf("error%v", errortype)}
	}
}

// ConnectToMessengerService connects a user to Online channel using a token
func (ms *MessengerService) connectToMessengerService(conn *socket.Socket, token string) bool {
	var err error
	var user *user.User

	MS, err := messengermanager.NewMessengerManager(nil)
	if err == nil {
		token := fmt.Sprintf("%v", token)
		user, err = MS.HasTokenAccess(token)
		if err == nil {
			AESkey := fmt.Sprintf("%v", conn.Data())

			// Removing socket past socket
			if ms.sockets[user.GetSocket()] != nil && user.GetSocket() != conn.Id() {
				ms.Logger.Printf("Disconnecting %s client", user.GetSocket())
				ms.sockets[user.GetSocket()].Disconnect(true)
				delete(ms.sockets, user.GetSocket())
			}

			user.SetSocketID(conn.Id())
			conn.SetData(gin.H{"key": AESkey, "user": *user})
			conn.Join("Online")
			ms.Logger.Printf("Sending User to %s", conn.Id())
			conn.Emit("Log In", user)
			return true
		}
	}
	ms.handleError(conn, "Login", err)
	return false
}

// sendMessage sends a message to group or chat
func (ms *MessengerService) sendMessage(conn *socket.Socket, Decryptedmessage string) (err error) {
	var decryptedContent string
	var message map[string]any = make(map[string]any)
	var group *group.Group
	context := conn.Data().(gin.H)

	//Decrypting content
	decryptedContent, err = utils.DecryptText(Decryptedmessage, context["key"].(string))

	json.Unmarshal([]byte(decryptedContent), &message)

	if err == nil {

		MS, err1 := messengermanager.NewMessengerManager(nil)
		err = err1
		if err == nil {
			toList := make([]*user.User, 0)

			err = mapstructure.Decode(message["to"], &toList)
			if err == nil {

				if len(toList) > 0 {

					if err == nil {
						var groupID primitive.ObjectID
						var sockets map[socket.SocketId]bool
						user := context["user"].(user.User)

						newMessage := msmessage.NewMessage(&user, message["content"].(string))

						// Checking if group
						groupID, err = MS.CheckGroup(user, toList)
						if err != nil {
							groupID, err = MS.CreateGroupByUsers(user, toList)
							if err == nil {

								group, err = MS.GetGroup(groupID)

								socketsIDS := MS.MapUsersToSocketsID(group.Members)

								if err == nil { // sending new Group to all members
									go func() {
										for socket := range socketsIDS {
											ms.NotifyChan <- &GeneralNotification{soketID: socket, data: group, NotificationType: "NewGroup"}
										}
									}()
								} else {
									ms.handleError(conn, "", err)
									return
								}

							}
						}

						newMessage.GroupID = groupID
						sockets, err = MS.SaveMessage(&user, toList, newMessage)

						if err == nil {

							go func() {

								for key := range sockets {
									if conn.Id() != key {
										ms.MessageSender <- &SocketMessage{message: newMessage, socket: key, messageType: "NewMessage"}
									}
								}

							}()

							conn.Emit("SentMessage", map[string]any{"ok": true, "message": gin.H{"id": newMessage.ID, "groupid": groupID}})

						} else {

							conn.Emit("SentMessage", gin.H{"ok": false, "message": gin.H{"groupid": newMessage.GroupID}})

						}

					}
				} else {
					err = errors.New("message need almost one user data")
				}
			}
		}
	}
	ms.handleError(conn, "", err)
	return
}

// getGroupHistory returns a list of 10 last messages using a date
func (ms *MessengerService) getGroupHistory(conn *socket.Socket, groupInfo map[string]any) (err error) {
	context := conn.Data().(gin.H)
	var ID primitive.ObjectID
	var history []*msmessage.Message
	var mTime time.Time
	var encryptedmessage string

	MS, err1 := messengermanager.NewMessengerManager(nil)
	err = err1
	if err == nil {
		ID, err = primitive.ObjectIDFromHex(groupInfo["ID"].(string))
		if err == nil {
			mTime, err = time.Parse(time.RFC3339, groupInfo["time"].(string))
			if err == nil {
				history, err = MS.GetGroupHistory(ID, mTime)
				if err == nil {
					for _, msg := range history {
						user := context["user"].(user.User)
						msg.WillSendtoUser(&user)
					}
					// re-using variable to encrypt
					encryptedmessage, err = utils.EncryptInterface(history, context["key"].(string))
					if err == nil {
						conn.Emit("History", encryptedmessage)
					}
				}
			}
		}
	}

	return
}

// seenMessage mark as Read a message by this connection user
func (ms *MessengerService) seenMessage(conn *socket.Socket, id string) (err error) {
	context := conn.Data().(gin.H)
	var message msmessage.Message
	var localUser user.User = context["user"].(user.User)
	var ID primitive.ObjectID

	ID, err = primitive.ObjectIDFromHex(id)

	if err == nil {

		MS, err := messengermanager.NewMessengerManager(nil)

		if err == nil {
			message, err = MS.MessageWasSeenBy(ID, localUser)
			if err == nil {
				if !message.From.IsEqual(&localUser) {
					message.WillSendtoUser(message.From)
					socket := MS.MapNumberToSocketID(message.From)
					if socket != nil {
						ms.MessageSender <- &SocketMessage{socket: *socket, message: &message, messageType: "ReadMessage"}
					}
				}

			}
		}
	}
	return
}

// GetGroups returns all user's group
func (ms *MessengerService) getGroups(conn *socket.Socket) (err error) {
	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok {
		user, _ := user.NewUser(contextMap["user"].(user.User))
		MS, err1 := messengermanager.NewMessengerManager(nil)
		err = err1
		if err == nil {
			var groups []*group.Group
			groups, err = MS.GetAllGroups(*user)

			conn.Emit("AllCurrentGroups", groups)
		}

	}
	ms.handleError(conn, "", err)
	return
}

// seenMessage marks as Read a message by this connection user
func (ms *MessengerService) createGroup(conn *socket.Socket, maps map[string]any) (err error) {
	var ingroup group.Group
	var newgroup *group.Group
	var jsonData []byte
	jsonData, err = json.Marshal(&maps)
	if err == nil {
		// Convert the JSON to a struct
		err = json.Unmarshal(jsonData, &ingroup)
		if err == nil {
			MS, err := messengermanager.NewMessengerManager(nil)

			if err == nil {
				ingroup.Admins = append(make([]*user.User, 0), ingroup.Members[0])

				newgroup, err = MS.CreateGroup(&ingroup)
				if err == nil {
					socketsIDS := MS.MapUsersToSocketsID(newgroup.Members)

					if err == nil { // sending new Group to all members
						go func() {
							for socket := range socketsIDS {
								ms.NotifyChan <- &GeneralNotification{soketID: socket, data: newgroup, NotificationType: "NewGroup"}
							}
						}()
					}
				}
			}
		}

	}

	return
}
