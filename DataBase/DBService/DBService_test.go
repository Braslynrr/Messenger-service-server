package dbservice_test

import (
	"MessengerService/dbservice"
	"MessengerService/user"
	"testing"
)

// TestNewDBService tests NewDBService
func TestNewDBService(t *testing.T) {
	_, err := dbservice.NewDBService()
	if err != nil {
		t.Logf("Instances should be initialized. Error: %v", err)
	}

}

func TestGetUser(t *testing.T) {
	var err error
	DS, _ := dbservice.NewDBService()
	temp := &user.User{Zone: "+506", Number: "00000000"}
	temp, err = DS.GetUser(*temp)
	if err != nil {
		t.Fatalf("Request should be done correctly. Error: %v", err)
	}
}
