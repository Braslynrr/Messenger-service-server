package main

import (
	"MessengerService/dbservice"
	messengerserviceapi "MessengerService/messengerserviceApi"
	"MessengerService/userserviceapi"
	"MessengerService/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func SetupServer(IsEncrypted bool) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	client := dbservice.MongoClient()
	store := mongodriver.NewStore(client, 3600, true, []byte("secret"))
	router.Use(sessions.Sessions("key", store))
	router.Use(messengerserviceapi.SetContextID)

	router.GET("/Key", messengerserviceapi.GetKey)
	router.POST("/User", utils.DecryptMiddleWare(IsEncrypted), userserviceapi.NewUser, utils.EncryptMiddleWare(IsEncrypted))
	router.POST("/User/Login", utils.DecryptMiddleWare(IsEncrypted), userserviceapi.Login, utils.EncryptMiddleWare(IsEncrypted))
	router.GET("/User/:zone/:number", userserviceapi.GetUser, utils.EncryptMiddleWare(IsEncrypted))
	router.GET("/Groups/:zone/:number", userserviceapi.GetGroups, utils.EncryptMiddleWare(IsEncrypted))

	handler := messengerserviceapi.NewSocketIo()

	router.GET("/socket.io/*any", gin.WrapH(handler))
	router.POST("/socket.io/*any", gin.WrapH(handler))

	router.StaticFS("/static", http.Dir("../../ServerFiles"))

	router.GET("/", messengerserviceapi.GetPage)

	return router
}

func main() {
	router := SetupServer(false)
	router.Run(":8080")
}
