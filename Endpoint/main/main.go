package main

import (
	messengerserviceapi "MessengerService/messengerserviceApi"
	"MessengerService/userserviceapi"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	socketioServer := messengerserviceapi.NewSocketIo()
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(ginsession.New())

	router.LoadHTMLFiles("websockets.html")
	router.GET("/", messengerserviceapi.GetPage)
	router.GET("/Key", messengerserviceapi.GetKey)
	router.POST("/User", userserviceapi.NewUser)
	router.POST("/User/Login", userserviceapi.Login)

	go socketioServer.Serve()
	defer socketioServer.Close()

	router.GET("/socket.io/*any", gin.WrapH(socketioServer))
	router.POST("/socket.io/*any", gin.WrapH(socketioServer))
	router.Run(":8080")
}
