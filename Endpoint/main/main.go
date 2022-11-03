package main

import (
	messengerserviceapi "MessengerService/messengerserviceApi"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(ginsession.New())

	router.GET("/MessengerService", messengerserviceapi.ConnectToMessengerService)
	router.LoadHTMLFiles("websockets.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "websockets.html", nil)
	})

	router.Run("localhost:8080")
}
