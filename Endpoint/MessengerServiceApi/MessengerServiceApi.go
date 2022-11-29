package messengerserviceapi

import (
	"MessengerService/messengermanager"
	"MessengerService/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	engineiopooling "github.com/googollee/go-socket.io/engineio/transport/polling"
)

const standarKeySize = 16

type Message struct {
	Info any `json:"info"`
}

// NewSocketIo creates a new Socketio server instance
func NewSocketIo() *socketio.Server {

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&engineiopooling.Transport{
				Client: &http.Client{
					Timeout: time.Hour,
				},
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		key, err := utils.GenerateRandomAESKey(standarKeySize)
		HandleError(s, "", err)
		s.SetContext(key)
		log.Println("Connected:", s.ID(), s.Namespace())

		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println(s.ID(), "closed due to", reason)
	})

	server.OnEvent("/", "messenger", connectToMessengerService)

	server.OnEvent("/", "sendMessage", HasUserMiddleware(sendMessage))
	return server
}

// HandleError handles erros using a type and error
func HandleError(conn socketio.Conn, errortype string, err error) {
	if err != nil {
		conn.Emit(fmt.Sprintf("error%v", errortype), gin.H{"error": err.Error()})
	}
}

// ConnectToMessengerService connects a user to Online channel using a token
func connectToMessengerService(conn socketio.Conn, args Message) {
	var err error
	var encryptedUser string
	MS, err := messengermanager.NewMessengerManager()
	token := fmt.Sprintf("%v", args.Info)
	user, err := MS.HasTokenAccess(token)
	if err == nil {
		AESkey := fmt.Sprintf("%v", conn.Context())
		user.SetSocketID(conn.ID())
		conn.SetContext(gin.H{"key": AESkey, "user": *user})
		conn.Join("Online")
		conn.Emit("WSKey", AESkey)
		UserInText, err := json.MarshalIndent(user, "", "")
		encryptedUser, err = utils.EncryptText(string(UserInText), AESkey)
		if err != nil {
			HandleError(conn, "", err)
			return
		}
		conn.Emit("Log In", encryptedUser)
		// return all user information

		return
	}
	HandleError(conn, "", err)
}

// sendMessage send message to group or chat
func sendMessage(conn socketio.Conn, message map[string]any) (err error) {
	// MS, err := messengermanager.NewMessengerManager()
	if err != nil {

	}
	HandleError(conn, "", err)
	return
}

// GetPage Return the Test page
func GetPage(c *gin.Context) {
	store := ginsession.FromContext(c)
	encryptKey, keyError := utils.GenerateRandomAESKey(16)
	if keyError == nil {
		store.Set("key", encryptKey)
		store.Save()
	}
	c.HTML(200, "websockets.html", nil)

}

// returns a key generated on GetPage
func GetKey(c *gin.Context) {
	store := ginsession.FromContext(c)
	data, ok := store.Get("key")
	if !ok {
		c.AbortWithStatus(500)
	}
	c.JSON(200, map[string]interface{}{"initialValue": data})
}

// HasUserMiddleware checks if user is loged in
func HasUserMiddleware(next func(socketio.Conn, map[string]any) error) func(socketio.Conn, map[string]any) error {
	return func(conn socketio.Conn, arg map[string]any) error {

		context := conn.Context()
		contextMap, ok := context.(gin.H)
		if ok {

			if contextMap["user"] != nil {
				return next(conn, arg)
			}

		}
		HandleError(conn, "", errors.New("connection has done an invalid action due to should log in."))
		return nil
	}
}
