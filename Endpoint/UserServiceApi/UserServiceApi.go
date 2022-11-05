package userserviceapi

import (
	"MessengerService/messengermanager"
	"MessengerService/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
