package messengerserviceapi

import (
	"MessengerService/dbservice"
	messengermanager "MessengerService/mesermanager"
	msmessage "MessengerService/message"
	"MessengerService/user"
	"MessengerService/userserviceapi"
	"MessengerService/utils"
	"errors"
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

func (ms *MessengerService) GetMessages(ctx *gin.Context) {
	var ok bool
	var err error
	var key string
	var body map[string]any = make(map[string]any)
	var ID primitive.ObjectID
	var history []*msmessage.Message
	var mTime time.Time
	var socketID socket.SocketId
	var encryptedmessage string

	if err = ctx.BindJSON(&body); err != nil {
		ctx.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	var temp string
	if temp, ok = body["socketID"].(string); !ok {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("there's no a valid authorized code"))
		return
	} else {
		socketID = socket.SocketId(temp)
		if ms.sockets[socketID] == nil && !ms.sockets[socketID].Connected() {
			ctx.AbortWithError(http.StatusUnauthorized, errors.New("there's no a valid authorized code"))
			return
		}
	}

	MS, err1 := messengermanager.NewMessengerManager(nil)
	err = err1
	if err == nil {
		ID, err = primitive.ObjectIDFromHex(body["ID"].(string))
		if err == nil {
			mTime, err = time.Parse(time.RFC3339, body["time"].(string))
			if err == nil {
				history, err = MS.GetGroupHistory(ID, mTime)
				if err == nil {
					user, ok := ms.sockets[socketID].Data().(gin.H)["user"].(user.User)
					if ok {
						key = ms.sockets[socketID].Data().(gin.H)["key"].(string)
						for _, msg := range history {

							msg.WillSendtoUser(&user)
						}

						// re-using variable to encrypt
						if len(history) != 0 {
							encryptedmessage, err = utils.EncryptInterface(history, key)
						} else {
							encryptedmessage = ""
						}

						if err == nil {

							ctx.String(http.StatusOK, encryptedmessage)
						}
					} else {
						ctx.AbortWithError(http.StatusUnauthorized, errors.New("you should be log in messenger to request messages"))
					}

				}
			}
		}
	}

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
		router.POST("Groups/Messages", utils.DecryptMiddleWare(IsEncrypted), ms.GetMessages, utils.EncryptMiddleWare(IsEncrypted))
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

var socketLock sync.Mutex

// MessageAndNotificationsnSender send some messages and notifications in background
func (ms *MessengerService) MessageAndNotificationsnSender() {
	for {
		select {
		case msg := <-ms.Sender:
			if ms.sockets[msg.socket] != nil && ms.sockets[msg.socket].Connected() {
				usercontext := ms.sockets[msg.socket].Data().(gin.H)
				encyptedMessage, err := utils.EncryptInterface(msg.message, usercontext["key"].(string))
				if err == nil {
					ms.Logger.Printf("Sending %s to %v", msg.messageType, msg.socket)
					socketLock.Lock()
					ms.sockets[msg.socket].Emit(msg.messageType, encyptedMessage)
					socketLock.Unlock()
				}
			}
		case err := <-ms.ErrorChan:
			ms.ErrorLog.Println("Error:", err.err.Error(), " from ", err.socket.Id())
			if err.socket.Connected() {
				err.socket.Emit(err.errorType, gin.H{"Error": err.err.Error()})
			}
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
