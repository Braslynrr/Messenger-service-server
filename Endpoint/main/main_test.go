package main

import (
	"MessengerService/dbservice"
	messengerserviceapi "MessengerService/messengerserviceApi"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/go-playground/assert"
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
	router, _ := ms.SetupServer(false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
