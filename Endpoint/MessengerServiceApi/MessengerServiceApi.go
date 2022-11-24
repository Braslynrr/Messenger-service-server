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
	"github.com/googollee/go-socket.io/engineio/transport/polling"
)

type Message struct {
	Info any `json:"info"`
}

func NewSocketIo() *socketio.Server {

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				Client: &http.Client{
					Timeout: time.Minute,
				}}}})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("Connected:", s.ID())
		s.Join("waiting")
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println(s.ID(), "closed due to", reason)
	})

	server.OnEvent("/", "messenger", ConnectToMessengerService)

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	return server
}

func ConnectToMessengerService(conn socketio.Conn, args string) error {
	fmt.Println("Aqui ando pa")
	MS, err := messengermanager.NewMessengerManager()
	token := fmt.Sprintf("%v", args)
	if user, err := MS.HasTokenAccess(token); err == nil {
		user.SetSocketID(conn.ID())
		conn.SetContext(*user)
		conn.Join("Online")
	}
	return err
}

func GetPage(c *gin.Context) {
	store := ginsession.FromContext(c)
	encryptKey, keyError := utils.GenerateRandomAESKey()
	if keyError == nil {
		store.Set("key", encryptKey)
		store.Save()
	}
	c.HTML(200, "websockets.html", nil)

}

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
