package messengerserviceapi

import (
	"MessengerService/dbservice"
	"MessengerService/group"
	messengermanager "MessengerService/mesermanager"
	msmessage "MessengerService/message"
	"MessengerService/user"
	"MessengerService/userserviceapi"
	"MessengerService/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocketMessage struct {
	message     *msmessage.Message
	socket      socket.SocketId
	messageType string
}
type SocketError struct {
	err       error
	errorType string
	socket    *socket.Socket
}

type MessengerService struct {
	Sender    chan *SocketMessage
	ErrorChan chan SocketError
	DoneChan  chan bool
	Wait      *sync.WaitGroup
	sockets   map[socket.SocketId]*socket.Socket
	Logger    *log.Logger
	ErrorLog  *log.Logger
	socketIO  *socket.Server
	DbService dbservice.DbInterface
	Sesion    gin.HandlerFunc
}

func (ms *MessengerService) newSocketIo() http.Handler {
	ms.socketIO = socket.NewServer(nil, nil)

	ms.socketIO.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)

		ms.Logger.Printf("New Connection: %v", client.Id())

		client.On("messenger", func(data ...any) {
			var err error

			if err == nil {

				WSKey, _ := utils.GenerateRandomAESKey(16)

				client.Emit("WSKey", WSKey)
				client.SetData(WSKey)

				token, ok := data[0].(string)
				if ok {
					ms.connectToMessengerService(client, token)
					ms.sockets[client.Id()] = client
				}
			} else {
				err = errors.New("object sent is not a token")
			}
			ms.handleError(client, "", err)
		})

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

		client.On("disconnect", func(reason ...any) {
			log.Print(reason)
		})

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
func (ms *MessengerService) connectToMessengerService(conn *socket.Socket, token string) {
	var err error
	var user *user.User

	MS, err := messengermanager.NewMessengerManager(nil)
	if err == nil {
		token := fmt.Sprintf("%v", token)
		user, err = MS.HasTokenAccess(token)
		if err == nil {
			AESkey := fmt.Sprintf("%v", conn.Data())
			user.SetSocketID(conn.Id())
			conn.SetData(gin.H{"key": AESkey, "user": *user})
			conn.Join("Online")
			conn.Emit("Log In", user)
			return
		}
	}
	ms.handleError(conn, "", err)
}

// sendMessage sends a message to group or chat
func (ms *MessengerService) sendMessage(conn *socket.Socket, Decryptedmessage string) (err error) {
	var decryptedContent string
	var message map[string]any = make(map[string]any)

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

								conn.Emit("NewGroup", group)

							}
						}

						newMessage.GroupID = groupID
						sockets, err = MS.SaveMessage(&user, toList, newMessage)

						if err == nil {

							go func() {

								for key := range sockets {
									if conn.Id() != key {
										ms.Sender <- &SocketMessage{message: newMessage, socket: key, messageType: "NewMessage"}
									}
								}

							}()

							encyptedMessage, err := utils.EncryptInterface(map[string]any{"ok": true, "message": newMessage}, context["key"].(string))
							if err == nil {
								conn.Emit("SentMessage", encyptedMessage)
							} else {
								conn.Emit("error", gin.H{"error": err.Error()})
							}

						} else {
							encryptedInterface, err = utils.EncryptInterface(newMessage, context["key"].(string))
							if err == nil {
								conn.Emit("SentMessage", gin.H{"ok": false, "message": encryptedInterface})
							}
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
						ms.Sender <- &SocketMessage{socket: *socket, message: &message, messageType: "ReadMessage"}
					}
				}

			}
		}
	}
	return
}

// GetPage Return the Test page
func (ms *MessengerService) getCodes(c *gin.Context) {
	c.File("../../ServerFiles/countrycodes/CountryCodes.json")
}

// returns a key generated on GetPage
func (ms *MessengerService) getKey(c *gin.Context) {

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

// hasUserMiddleware checks if user is loged in
func hasUserMiddleware[T any](conn *socket.Socket, next func(*socket.Socket, T) error, args ...any) (err error) {

	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok {
		arg, ok := args[0].(T)
		if ok && contextMap["user"] != nil {

			return next(conn, arg)
		}

	}

	return errors.New("connection should be connected using a token")
}

// HasUserMiddleware checks if user is loged in
func hasUserMiddlewareNoParam(conn *socket.Socket, next func(*socket.Socket) error) (err error) {

	context := conn.Data()
	contextMap, ok := context.(gin.H)
	if ok && contextMap["user"] != nil {

		return next(conn)
	}

	return errors.New("connection should be connected using a token")
}

// SetContextID set token for session
func setContextID(ctx *gin.Context) { // pending to fix

	sessionID, err := ctx.Cookie("session_id")

	if err != nil || sessionID == "" {
		sessionID, _ = utils.GenerateToken()
		ctx.SetCookie("session_id", sessionID, 0, "", "", false, true)

		session := sessions.Default(ctx)
		session.Set("session_id", sessionID)
		session.Save()

		defer session.Save()
	}

	ctx.Next()
}

// SetupServer sets up the Gin server
func (ms *MessengerService) SetupServer(IsEncrypted bool) (router *gin.Engine, err error) {

	mm, err := messengermanager.NewMessengerManager(ms.DbService)
	if mm.IsInitialize() {
		ms.sockets = make(map[socket.SocketId]*socket.Socket)

		router = gin.Default()
		router.Use(cors.Default())
		router.Use(ms.Sesion)
		router.Use(setContextID)

		router.GET("/Key", ms.getKey)
		router.POST("/User", utils.DecryptMiddleWare(IsEncrypted), userserviceapi.NewUser, utils.EncryptMiddleWare(IsEncrypted))
		router.POST("/User/Login", utils.DecryptMiddleWare(IsEncrypted), userserviceapi.Login, utils.EncryptMiddleWare(IsEncrypted))
		router.GET("/User/:zone/:number", userserviceapi.GetUser, utils.EncryptMiddleWare(IsEncrypted))
		router.GET("/Groups/:zone/:number", userserviceapi.GetGroups, utils.EncryptMiddleWare(IsEncrypted))
		router.GET("/CountryCodes", ms.getCodes)

		handler := ms.newSocketIo()

		router.GET("/socket.io/*any", gin.WrapH(handler))
		router.POST("/socket.io/*any", gin.WrapH(handler))

		router.Use(static.Serve("/", static.LocalFile("../../ServerFiles/messenger-ui/build", true)))
	} else {
		err = errors.New("messenger manager is not initialized")
	}
	return
}

// MessageAndNotificationsnSender send some messages and notifications in background
func (ms *MessengerService) MessageAndNotificationsnSender() {
	for {
		select {
		case msg := <-ms.Sender:
			usercontext := ms.sockets[msg.socket].Data().(gin.H)
			encyptedMessage, err := utils.EncryptInterface(msg.message, usercontext["key"].(string))
			if err == nil {
				ms.Logger.Println("Sending ", msg.messageType, " to ", msg.socket)
				ms.sockets[msg.socket].Emit(msg.messageType, encyptedMessage)
			}
		case err := <-ms.ErrorChan:
			ms.ErrorLog.Println("Error:", err.err.Error(), " from ", err.socket.Id())
			err.socket.Emit(err.errorType, gin.H{"Error": err.err.Error()})
		case <-ms.DoneChan:
			return
		}
	}
}

// ListenForShutdown Listens for a signal to shutdown
func (ms *MessengerService) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ms.ShutDown()
	os.Exit(0)
}

// ShutDown shutdowns the server
func (ms *MessengerService) ShutDown() {
	ms.Logger.Println("Would run cleanup tasks...")

	ms.Logger.Println("Closign SocketIo Server...")

	ms.socketIO.To("Online").DisconnectSockets(true)
	ms.socketIO.Clear()
	ms.socketIO = nil

	ms.Wait.Wait()

	ms.Logger.Println("Closign Channels and shutting down the server...")
	close(ms.Sender)
	close(ms.ErrorChan)
	close(ms.DoneChan)
	ms.sockets = nil

}
