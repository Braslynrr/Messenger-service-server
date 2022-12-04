package main

import (
	messengerserviceapi "MessengerService/messengerserviceApi"
	"MessengerService/userserviceapi"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	store := cookie.NewStore([]byte(""))
	router.Use(sessions.Sessions("key", store))

	router.GET("/Key", messengerserviceapi.GetKey)
	router.POST("/User", userserviceapi.NewUser)
	router.POST("/User/Login", userserviceapi.Login)

	handler := messengerserviceapi.NewSocketIo()

	router.GET("/socket.io/*any", gin.WrapH(handler))
	router.POST("/socket.io/*any", gin.WrapH(handler))
	router.LoadHTMLFiles("../../ServerFiles/html/websockets.html")

	router.GET("/", messengerserviceapi.GetPage)

	router.Run(":8080")
}
