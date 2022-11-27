package messengerserviceapi

import (
	"MessengerService/messengermanager"
	"MessengerService/utils"
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

	server.OnEvent("/", "messenger", ConnectToMessengerService)

	return server
}

// HandleError handles erros using a type and error
func HandleError(conn socketio.Conn, errortype string, err error) {
	if err != nil {
		conn.Emit(fmt.Sprintf("error%v", errortype), gin.H{"error": err.Error()})
	}
}

// ConnectToMessengerService connects a user to Online channel using a token
func ConnectToMessengerService(conn socketio.Conn, args Message) {
	MS, err := messengermanager.NewMessengerManager()
	token := fmt.Sprintf("%v", args.Info)
	user, err := MS.HasTokenAccess(token)
	if err == nil {
		AESkey := fmt.Sprintf("%v", conn.Context())
		//user.Password, err = utils.EncryptText(AESkey, user.Password)
		user.SetSocketID(conn.ID())
		conn.SetContext(gin.H{"key": AESkey, "user": *user})
		conn.Join("Online")
		conn.Emit("Log In", user)
		// return all user information
		return
	}
	HandleError(conn, "", err)
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

func SocketServer(c *gin.Context) {

	NewSocketIo().ServeHTTP(c.Writer, c.Request)

}
