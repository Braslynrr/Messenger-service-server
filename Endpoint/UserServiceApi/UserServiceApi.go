package userserviceapi

import (
	"MessengerService/messengermanager"
	"MessengerService/user"
	"MessengerService/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/gorilla/websocket"
)

// NewUser recieves post request to regsister a new user
func NewUser(c *gin.Context) {
	tempUser := &user.User{}
	var encryptedUser string
	var mapUser map[string]string
	var err error

	if err = c.BindJSON(&mapUser); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	store := ginsession.FromContext(c)
	key, _ := store.Get("key")

	encryptedUser, err = utils.DecryptText(mapUser["user"], fmt.Sprint(key))

	err = json.Unmarshal([]byte(encryptedUser), tempUser)

	if err == nil {

		messman, err2 := messengermanager.NewMessengerManager()
		err = err2
		if err2 == nil {
			var ok bool
			ok, err = messman.InsertUser(*tempUser)

			if ok && err == nil {

				c.Done()
				return
			}
		}
	}
	c.IndentedJSON(http.StatusNotAcceptable, err)
}

// Login checks if user exists to returns a new token to connect
func Login(c *gin.Context) {
	tempUser := &user.User{}
	var encryptedUser string
	var mapUser map[string]string
	var err error
	var token string

	if err = c.BindJSON(&mapUser); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}
	store := ginsession.FromContext(c)
	key, _ := store.Get("key")

	encryptedUser, err = utils.DecryptText(mapUser["user"], fmt.Sprint(key))

	err = json.Unmarshal([]byte(encryptedUser), tempUser)

	if err == nil {

		tempUser, err = user.NewUser(*tempUser)

		if err == nil {

			messman, mmerr := messengermanager.NewMessengerManager()

			if mmerr == nil {

				if err == nil {

					token, err = messman.Login(*tempUser)

					if err == nil {

						c.IndentedJSON(http.StatusOK, gin.H{"token": token})
						return
					}
				}

			} else {
				err = mmerr
			}
		}
	}
	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
}

// ConnectUser verifies and connects an user with a token
func ConnectUser(token string, conn *websocket.Conn, encryptKey string) error {
	MM, err := messengermanager.NewMessengerManager()
	if err == nil {
		err = MM.ConnectUser(token, conn)
	}
	return err
}
