package main

import (
	"MessengerService/dbservice"
	messengerserviceapi "MessengerService/messengerserviceApi"
	"log"
	"os"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
)

func main() {
	var DB *dbservice.DbService
	var router *gin.Engine
	var err error
	DB, err = dbservice.NewDBService()
	logger := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	logerror := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err == nil {

		client := dbservice.MongoClient()
		store := mongodriver.NewStore(client, 3600, true, []byte("secret"))

		ms := messengerserviceapi.MessengerService{
			Sender:    make(chan *messengerserviceapi.SocketMessage, 100),
			ErrorChan: make(chan messengerserviceapi.SocketError, 100),
			DoneChan:  make(chan bool),
			Wait:      &sync.WaitGroup{},
			Logger:    logger,
			ErrorLog:  logerror,
			DbService: DB,
			Sesion:    sessions.Sessions("key", store),
		}

		router, err = ms.SetupServer(false)

		//listen messages and notifications
		if err == nil {
			go ms.MessageAndNotificationsnSender()

			//listen for shutdown
			go ms.ListenForShutdown()

			router.Run(":8080")
		}

	} else {
		logerror.Printf(err.Error())
	}
}
