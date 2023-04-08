package main

import (
	"MessengerService/dbservice"
	messengerserviceapi "MessengerService/messengerserviceApi"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/go-playground/assert"
)

func TestMain(t *testing.T) {
	ms := messengerserviceapi.MessengerService{
		Sender:    make(chan *messengerserviceapi.SocketMessage, 1),
		ErrorChan: make(chan messengerserviceapi.SocketError, 1),
		DoneChan:  make(chan bool),
		Wait:      &sync.WaitGroup{},
		Logger:    log.Default(),
		DbService: &dbservice.DbTest{},
	}
	router, _ := ms.SetupServer(false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	fi, _ := os.Open("../../ServerFiles/html/websockets.html")
	buf := make([]byte, w.Body.Len())
	fi.Read(buf)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, buf, w.Body.Bytes())
}
