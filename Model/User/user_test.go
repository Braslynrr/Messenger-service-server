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

// TestNewUserNumberFail calls NewUser checking a Number property raise a exception
func TestNewUserNumberFail(t *testing.T) {
	user1, err := user.NewUser(user.User{Zone: "+000", Number: "00000l00", Password: "123"})
	if user1 != nil {
		t.Fatalf("User should be null. User: %v", user1)
	}

	if err == nil {
		t.Fatalf("An Exception should be returned")
	}
}

// TestNewUserZoneFail calls NewUser checking a Zone property raise a exception
func TestNewUserZoneFail(t *testing.T) {
	user1, err := user.NewUser(user.User{Zone: "+0l0", Number: "00000000", Password: "123"})
	if user1 != nil {
		t.Fatalf("User should be null. User: %v", user1)
	}

	if err == nil {
		t.Fatalf("An Exception should be returned")
	}
}

// TestNewUserPasswordFail calls NewUser checking a Password property raise a exception
func TestNewUserPasswordFail(t *testing.T) {
	user1, err := user.NewUser(user.User{Zone: "+000", Number: "00000000"})
	if user1 != nil {
		t.Fatalf("User should be null. User: %v", user1)
	}

	if err == nil {
		t.Fatalf("An Exception should be returned")
	}
}

// TestIsEqual calls IsEqual checking two users are equal
func TestIsEqual(t *testing.T) {
	user1, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	user2, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})

	if !user1.IsEqual(user2) {
		t.Fatalf("Both users should be equal. %v1 != %v2", user1, user2)
	}

	user2.UserName = "Antonio"
	user2.Password = "paco"
	if !user1.IsEqual(user2) {
		t.Fatalf("Both modified users should be equal. %v1 != %v2", user1, user2)
	}

	user2.Zone = "+001"
	if user1.IsEqual(user2) {
		t.Fatalf("Both users should be different. %v1 != %v2", user1, user2)
	}
}

// TestCredentials calls credentials checking equality.
func TestCredentials(t *testing.T) {
	user1, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	user2, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})

	if !user1.Credentials(user2) {
		t.Fatalf("Password, zone and number from user1 and user2 should be alike. %v == %v", user1, user2)
	}
}
