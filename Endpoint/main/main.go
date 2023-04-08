package main

import (
	messengerserviceapi "MessengerService/messengerserviceApi"
	"log"
	"sync"
)

func main() {
	ms := messengerserviceapi.MessengerService{
		Sender:    make(chan *messengerserviceapi.SocketMessage, 100),
		ErrorChan: make(chan messengerserviceapi.SocketError, 100),
		DoneChan:  make(chan bool),
		Wait:      &sync.WaitGroup{},
		Logger:    log.Default(),
	}

	//listen messages and notifications
	go ms.MessageAndNotificationsnSender()

	//listen for shutdown
	go ms.ListenForShutdown()

	router := ms.SetupServer(false)
	router.Run(":8080")
}
