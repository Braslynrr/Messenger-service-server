package main

import (
	messengerserviceapi "MessengerService/messengerserviceApi"
	"MessengerService/userserviceapi"

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
	router.GET("/", messengerserviceapi.GetPage)
	router.GET("/Key", messengerserviceapi.GetKey)
	router.POST("/User", userserviceapi.NewUser)

	router.Run("localhost:8080")
}
