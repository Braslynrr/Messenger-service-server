package main

import (
	"MessengerService/dbservice"
	messengerserviceapi "MessengerService/messengerserviceApi"
	"log"
	"sync"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func TestMain(t *testing.T) {
	store := cookie.NewStore([]byte(""))

	ms := messengerserviceapi.MessengerService{
		MessageSender:   make(chan *messengerserviceapi.SocketMessage, 100),
		MessageDoneChan: make(chan bool),
		ErrorChan:       make(chan messengerserviceapi.SocketError, 100),
		ErrorDoneChan:   make(chan bool),
		NotifyChan:      make(chan *messengerserviceapi.GeneralNotification),
		NotifyDoneChan:  make(chan bool),
		Wait:            &sync.WaitGroup{},
		Logger:          log.Default(),
		DbService:       &dbservice.DbTest{},
		Sesion:          sessions.Sessions("key", store),
	}
	_, err := ms.SetupServer(false)

	if err != nil {
		t.Errorf("SetupServer should not fail. Err: %s", err.Error())
	}
}
