package messengerserviceapi

import (
	"MessengerService/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Action string `json:"action"`
	Info   any    `json:"info"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ConnectToMessengerService Conects the gin server with the WS server,
func ConnectToMessengerService(c *gin.Context) {

	MessengerServiceHandler(c.Writer, c.Request)
}

//printError prints error to the user
func printError(connection *websocket.Conn, err error) {
	connection.WriteJSON(map[string]interface{}{"action": "notify", "status": err})
}

func MessengerServiceHandler(responseWriter gin.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(responseWriter, request, nil)

	var connected bool = true
	if err == nil {
		encryptKey, keyError := utils.GenerateRandomAESKey()
		conn.WriteJSON(map[string]interface{}{"action": "UpdateEncryptKey", "Key": encryptKey})

		for connected && keyError == nil {
			// Read message from browser
			var message Message
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			jsonInfo := fmt.Sprintf("%s", msg)
			value := []byte(jsonInfo)
			json.Unmarshal(value, &message)
			message.Action = strings.ToLower(message.Action)
			switch message.Action {
			default:
				printError(conn, errors.New("Action does not exist."))
			}

		}

		if keyError != nil {
			printError(conn, errors.New("An Error ocurred while encrypting key was generating."))
		}

		conn.Close()
	}
}
