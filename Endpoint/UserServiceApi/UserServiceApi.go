package userserviceapi

import (
	"MessengerService/group"
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
	var bytes []byte
	var err error

	if bytes, err = c.GetRawData(); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}

	session := sessions.Default(c)
	key := session.Get("key")

	encryptedUser, err = utils.DecryptText(string(bytes), fmt.Sprint(key))
	if err == nil {

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
	}
	c.IndentedJSON(http.StatusNotAcceptable, err)
}

// Login checks if user exists to returns a new token to connect
func Login(c *gin.Context) {
	tempUser := &user.User{}
	var bytes []byte
	var encryptedUser string
	var err error
	var token string

	if bytes, err = c.GetRawData(); err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, err)
		return
	}
	session := sessions.Default(c)
	key := session.Get("key")

	encryptedUser, err = utils.DecryptText(string(bytes), fmt.Sprint(key))

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

// GetUser Gets an user by its zone and number
func GetUser(c *gin.Context) {
	var tempUser *user.User
	var err error

	if err = c.ShouldBindUri(&tempUser); err == nil {

		messman, mmerr := messengermanager.NewMessengerManager()

		if mmerr == nil {
			tempUser, err = messman.GetUser(*tempUser)
			if err == nil {
				c.IndentedJSON(http.StatusOK, tempUser)
				return
			}
		} else {
			err = mmerr
		}
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
}

// GetGroups Gets groups by user
func GetGroups(c *gin.Context) {
	var tempUser *user.User
	var grouplist []group.Group
	var err error

	if err = c.ShouldBindUri(&tempUser); err == nil {

		messman, mmerr := messengermanager.NewMessengerManager()

		if mmerr == nil {
			grouplist, err = messman.GetAllGroups(*tempUser)
			if err == nil {
				c.IndentedJSON(http.StatusOK, grouplist)
				return
			}
		} else {
			err = mmerr
		}
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
}
