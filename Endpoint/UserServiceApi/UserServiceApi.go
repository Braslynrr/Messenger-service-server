package userserviceapi

import (
	messengermanager "MessengerService/mesermanager"
	"MessengerService/user"
	"MessengerService/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

	session := sessions.Default(c)
	key := session.Get("key")

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
	session := sessions.Default(c)
	key := session.Get("key")

	encryptedUser, err = utils.DecryptText(mapUser["user"], fmt.Sprint(key))

	if err == nil {

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
	}
	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
}
