package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/assert"
)

func TestMain(t *testing.T) {
	router := SetupServer(false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	fi, _ := os.Open("../../ServerFiles/html/websockets.html")
	buf := make([]byte, w.Body.Len())
	fi.Read(buf)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, buf, w.Body.Bytes())
}
