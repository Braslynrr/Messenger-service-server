package user_test

import (
	"MessengerService/user"
	"testing"
)

// TestNewUser calls NewUser checking a new user is returned
func TestNewUser(t *testing.T) {
	user1 := user.NewUser(user.User{})
	if user1 == nil {
		t.Fatalf("User is null. User: %v", user1)
	}
}
