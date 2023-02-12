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
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const standarKeySize = 16

func NewSocketIo() http.Handler {
	serverIO := socket.NewServer(nil, nil)

	serverIO.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)

		log.Printf("New Connection: %v", client.Id())

		WSKey, err := utils.GenerateRandomAESKey(standarKeySize)
		if err == nil {
			client.SetData(WSKey)
			client.Emit("WSKey", WSKey)
		}

		HandleError(client, "", err)
		client.On("messenger", func(data ...any) {

			token, ok := data[0].(string)
			if ok {
				context, ok := client.Data().(gin.H)
				if ok && context["user"] != nil {
					sendUserInfo(client, nil)
				} else {
					connectToMessengerService(client, token)
				}
			} else {
				err = errors.New("object sended is not a token")
			}

		})

		client.On("sendMessage", func(message ...any) {

			HandleError(client, "", HasUserMiddleware(client, sendMessage, message...))
		})

		client.On("GetCurrentGroups", func(...any) {
			HandleError(client, "", HasUserMiddlewareNoParam(client, GetGroups))
		})

		client.On("GroupHistory", func(history ...any) {
			HandleError(client, "", HasUserMiddleware(client, getGroupHistory, history...))
		})

		client.On("SendSeen", func(SeenMesage ...any) {
			HandleError(client, "", HasUserMiddleware(client, seenMessage, SeenMesage...))
		})

		client.On("disconnect", func(reason ...any) {
			log.Print(reason)
		})
	})
	return serverIO.ServeHandler(nil)
}

// HandleError handles erros using a type and error
func HandleError(conn *socket.Socket, errortype string, err error) {
	if err != nil {
		conn.Emit(fmt.Sprintf("error%v", errortype), gin.H{"error": err.Error()})
	}
}

// ConnectToMessengerService connects a user to Online channel using a token
func connectToMessengerService(conn *socket.Socket, token string) {
	var err error

	MS, err := messengermanager.NewMessengerManager()
	if err == nil {
		token := fmt.Sprintf("%v", token)
		user, err := MS.HasTokenAccess(token)
		if err == nil {
			AESkey := fmt.Sprintf("%v", conn.Data())
			user.SetSocketID(conn.Id())
			conn.SetData(gin.H{"key": AESkey, "user": *user})
			conn.Join("Online")
			sendUserInfo(conn, user)
			return
		}
	}
	HandleError(conn, "", err)
}

func sendUserInfo(conn *socket.Socket, connUser *user.User) {
	var encryptedUser string
	var err error
	context, ok := conn.Data().(gin.H)
	if ok {
		AESkey := fmt.Sprintf("%v", context["key"])
		if connUser == nil {
			var CU user.User = context["user"].(user.User)
			connUser = &CU
		}
		encryptedUser, err = utils.EncryptInterface(connUser, AESkey)
		if err != nil {
			HandleError(conn, "", err)
			return
		}
	}

	conn.Emit("Log In", encryptedUser)
}

// sendMessage sends a message to group or chat
func sendMessage(conn *socket.Socket, Decryptedmessage string) (err error) {
	var decryptedContent string
	var message map[string]any = make(map[string]any)

	context := conn.Data().(gin.H)

	//Decrypting content
	decryptedContent, err = utils.DecryptText(Decryptedmessage, context["key"].(string))

	json.Unmarshal([]byte(decryptedContent), &message)

	if err == nil {

		MS, err1 := messengermanager.NewMessengerManager()
		err = err1
		if err == nil {
			toList := make([]*user.User, 0)

			err = mapstructure.Decode(message["to"], &toList)
			if err == nil {

				if len(toList) > 0 {

					if err == nil {

						var encryptedInterface string
						var groupID primitive.ObjectID
						var sockets map[socket.SocketId]bool
						user := context["user"].(user.User)

						newMessage := msmessage.NewMessage(&user, message["content"].(string))

						// Checking if group
						groupID, err = MS.CheckGroup(user, toList)
						if err != nil {
							groupID, err = MS.CreateGroup(user, toList)
							if err == nil {
								group, _ := MS.GetGroup(groupID)
								encryptedInterface, err = utils.EncryptInterface(group, context["key"].(string))
								if err == nil {
									conn.Emit("NewGroup", encryptedInterface)
								}
							}
						}

						newMessage.GroupID = groupID
						sockets, err = MS.SaveMessage(&user, toList, newMessage)

						if err == nil {
							MS.SendToNumber(conn, "NewMessage", sockets, newMessage)
							encyptedMessage, err := utils.EncryptInterface(map[string]any{"ok": true, "message": newMessage}, context["key"].(string))
							if err == nil {
								conn.Emit("SendedMessage", encyptedMessage)
							} else {
								conn.Emit("error", gin.H{"error": err.Error()})
							}

						} else {
							encryptedInterface, err = utils.EncryptInterface(newMessage, context["key"].(string))
							if err == nil {
								conn.Emit("SendedMessage", gin.H{"ok": false, "message": encryptedInterface})
							}
						}

					}
				} else {
					err = errors.New("message need almost one user data")
				}
			}
		}
	}
	HandleError(conn, "", err)
	return
}

// getGroupHistory returns a list of 10 last messages using a date
func getGroupHistory(conn *socket.Socket, decryptedmessage string) (err error) {
	context := conn.Data().(gin.H)
	var ID primitive.ObjectID
	var history []*msmessage.Message
	var mTime time.Time
	var info map[string]any

	decryptedmessage, err = utils.DecryptText(decryptedmessage, context["key"].(string))
	if err == nil {
		err = json.Unmarshal([]byte(decryptedmessage), &info)
		if err == nil {
			MS, err1 := messengermanager.NewMessengerManager()
			err = err1
			if err == nil {
				ID, err = primitive.ObjectIDFromHex(info["ID"].(string))
				if err == nil {
					mTime, err = time.Parse(time.RFC3339, info["time"].(string))
					if err == nil {
						history, err = MS.GetGroupHistory(ID, mTime)
						if err == nil {
							for _, msg := range history {
								user := context["user"].(user.User)
								msg.WillSendtoUser(&user)
							}
							// re-using variable to encrypt
							decryptedmessage, err = utils.EncryptInterface(history, context["key"].(string))
							if err == nil {
								conn.Emit("History", decryptedmessage)
							}
						}
					}
				}
			}
		}
	}

	return
}

// seenMessage mark as Read a message by this connection user
func seenMessage(conn *socket.Socket, decryptedmessage string) (err error) {
	context := conn.Data().(gin.H)
	var ID primitive.ObjectID
	var message msmessage.Message
	var localUser user.User = context["user"].(user.User)

	decryptedmessage, err = utils.DecryptText(decryptedmessage, context["key"].(string))
	if err == nil {
		MS, err1 := messengermanager.NewMessengerManager()
		err = err1
		if err == nil {
			ID, err = primitive.ObjectIDFromHex(decryptedmessage)
			if err == nil {
				message, err = MS.MessageWasSeenBy(ID, localUser)
				if err == nil {
					if !message.From.IsEqual(&localUser) {
						message.WillSendtoUser(message.From)
						socket := MS.MapNumberToSocketID(message.From)
						MS.SendToNumber(conn, "ReadedMessage", socket, &message)
					}

				}

			}
		}
	}
	return
}

// GetPage Return the Test page
func GetPage(c *gin.Context) {
	c.HTML(200, "websockets.html", nil)
}

// returns a key generated on GetPage
func GetKey(c *gin.Context) {

	var err error
	session := sessions.Default(c)
	encryptKey, err := utils.GenerateRandomAESKey(16)
	if err == nil {
		session.Set("key", encryptKey)
		err = session.Save()
	}
	if err != nil {
		c.AbortWithStatus(500)
	}
	c.JSON(200, map[string]interface{}{"initialValue": encryptKey})
}

// GetGroups returns all user's group
func GetGroups(conn *socket.Socket) (err error) {
	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok {
		user, _ := user.NewUser(contextMap["user"].(user.User))
		MS, err1 := messengermanager.NewMessengerManager()
		err = err1
		if err == nil {
			var groups []group.Group
			groups, err = MS.GetAllGroups(*user)
			if err == nil {
				var encryptedGroups string
				context := conn.Data().(gin.H)
				encryptedGroups, err = utils.EncryptInterface(groups, context["key"].(string))
				if err == nil {
					conn.Emit("AllCurrentGroups", encryptedGroups)
				}
			}

		}
	}
	HandleError(conn, "", err)
	return
}

// HasUserMiddleware checks if user is loged in
func HasUserMiddleware(conn *socket.Socket, next func(*socket.Socket, string) error, args ...any) (err error) {

	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok {
		arg, ok := args[0].(string)
		if ok && contextMap["user"] != nil {

			return next(conn, arg)
		}

	}

	return errors.New("connection has done an invalid action due to should log in")
}

// HasUserMiddleware checks if user is loged in
func HasUserMiddlewareNoParam(conn *socket.Socket, next func(*socket.Socket) error) (err error) {

	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok && contextMap["user"] != nil {

		return next(conn)
	}

	return errors.New("connection has done an invalid action due to should log in")
}
