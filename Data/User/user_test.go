package user_test

import (
	"MessengerService/user"
	"testing"
)

// TestNewUser calls NewUser checking a new user is returned
func TestNewUser(t *testing.T) {
	user1, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	if user1 == nil {
		t.Fatalf("User is null. User: %v", user1)
	}
}

//TestNewUserNumberFail calls NewUser checking a Number property raise a exception
func TestNewUserNumberFail(t *testing.T) {
	user1, err := user.NewUser(user.User{Zone: "+000", Number: "00000l00", Password: "123"})
	if user1 != nil {
		t.Fatalf("User should be null. User: %v", user1)
	}

	if err == nil {
		t.Fatalf("An Exception should be returned")
	}
}

//TestNewUserZoneFail calls NewUser checking a Zone property raise a exception
func TestNewUserZoneFail(t *testing.T) {
	user1, err := user.NewUser(user.User{Zone: "+0l0", Number: "00000000", Password: "123"})
	if user1 != nil {
		t.Fatalf("User should be null. User: %v", user1)
	}

	if err == nil {
		t.Fatalf("An Exception should be returned")
	}
}

//TestNewUserPasswordFail calls NewUser checking a Password property raise a exception
func TestNewUserPasswordFail(t *testing.T) {
	user1, err := user.NewUser(user.User{Zone: "+000", Number: "00000000"})
	if user1 != nil {
		t.Fatalf("User should be null. User: %v", user1)
	}

	if err == nil {
		t.Fatalf("An Exception should be returned")
	}
}
