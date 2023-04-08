package userserviceapi_test

import (
	"MessengerService/dbservice"
	messengermanager "MessengerService/mesermanager"
	"MessengerService/user"
	"MessengerService/userserviceapi"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var db dbservice.DbInterface = &dbservice.DbTest{}

func TestInsertUser(t *testing.T) {
	messengermanager.NewMessengerManager(db)

	temp := user.User{Zone: "+506", Number: "62073447", Password: "poncho"}

	r := gin.Default()
	r.POST("/NewUser", userserviceapi.NewUser)

	w := httptest.NewRecorder()
	data, _ := json.Marshal(temp)

	req, _ := http.NewRequest("POST", "/NewUser", strings.NewReader(string(data)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// TestLogin test a user can login
func TestLogin(t *testing.T) {
	messengermanager.NewMessengerManager(db)

	temp := user.User{Zone: "+506", Number: "62073447", Password: "poncho"}

	db.InsertUser(temp)

	r := gin.Default()
	r.POST("/login", userserviceapi.Login)

	w := httptest.NewRecorder()
	data, _ := json.Marshal(temp)

	req, _ := http.NewRequest("POST", "/login", strings.NewReader(string(data)))
	r.ServeHTTP(w, req)

	fmt.Println(w.Body)
	assert.Equal(t, 200, w.Code)
}

// TestGetUser test a user request from client
func TestGetUser(t *testing.T) {
	messengermanager.NewMessengerManager(db)

	var request *user.User
	temp := user.User{Zone: "+506", Number: "62073447"}
	db.InsertUser(temp)

	r := gin.Default()
	r.GET("/User/:zone/:number", userserviceapi.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/User/+506/62073447", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &request)

	assert.Equal(t, temp.IsEqual(request), true)
}
