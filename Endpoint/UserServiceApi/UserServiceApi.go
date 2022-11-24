package userserviceapi

import (
	"MessengerService/messengermanager"
	"MessengerService/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NewUser recieves post request to regsister a new user
func NewUser(c *gin.Context) {
	var tempUser *user.User

	if err := c.BindJSON(&tempUser); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
	}

	tempUser, err := user.NewUser(*tempUser)

	if err == nil {

		messman, err := messengermanager.NewMessengerManager()

		if err == nil {

			ok, err := messman.InsertUser(*tempUser)

			if ok && err == nil {

				c.Done()
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotAcceptable, err)
}

func Login(c *gin.Context) {
	var tempUser *user.User
	var err error
	var token string

	if err = c.BindJSON(&tempUser); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	tempUser, err = user.NewUser(*tempUser)

	if err == nil {

		messman, mmerr := messengermanager.NewMessengerManager()

		if mmerr == nil {

			token, err = messman.Login(*tempUser)

			if err == nil {

				c.IndentedJSON(http.StatusOK, gin.H{"token": token})
				return
			}
		} else {
			err = mmerr
		}
	}
	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
}

func ConnectUser(token string, conn *websocket.Conn, encryptKey string) error {
	MM, err := messengermanager.NewMessengerManager()
	if err == nil {
		err = MM.ConnectUser(token, conn)
	}
	return err
}

//err := json.Unmarshal([]byte(fmt.Sprintf("%v", info)), &user)
