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
	"github.com/joho/godotenv"
)

func main() {

	/*
		Just the file name when the main.go is in the root, if not, set the directory /local.env.
		for mor references go to https://github.com/joho/godotenv
	*/

	err := godotenv.Load("local.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var DB *dbservice.DbService
	var router *gin.Engine
	DB, err = dbservice.NewDBService()
	logger := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	logerror := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err == nil {

		client := dbservice.MongoClient()
		store := mongodriver.NewStore(client, 3600, true, []byte("secret"))

		ms := messengerserviceapi.MessengerService{
			MessageSender:   make(chan *messengerserviceapi.SocketMessage, 100),
			MessageDoneChan: make(chan bool),
			ErrorChan:       make(chan messengerserviceapi.SocketError, 100),
			ErrorDoneChan:   make(chan bool),
			NotifyChan:      make(chan *messengerserviceapi.GeneralNotification, 100),
			NotifyDoneChan:  make(chan bool),
			Wait:            &sync.WaitGroup{},
			Logger:          logger,
			ErrorLog:        logerror,
			DbService:       DB,
			Sesion:          sessions.Sessions("key", store),
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
