package group_test

import (
	"MessengerService/group"
	"MessengerService/user"
	"testing"
)

// TestNewGroupFail calls NewGroup checking an error
func TestNewGroupFail(t *testing.T) {
	user, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	_, err := group.NewGroup(user)
	if err == nil {
		t.Fatalf("TestNewGroup Should return a error. Err: %v", err)
	}
}

// TestNewGroup calls NewGroup checking a new group is correctly created
func TestNewGroup(t *testing.T) {
	user1, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	user2, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	newGroup, err := group.NewGroup(user1, user2)
	if err != nil {
		t.Fatalf("An error should not be returned. Error: %v", err)
	}
	if newGroup == nil {
		t.Fatalf("The new group should not be nil.")
	}
}

// TestMinimunUsers calls NewGroup checking the new group have at least2 members
func TestMinimunUsers(t *testing.T) {
	user1, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	user2, _ := user.NewUser(user.User{Zone: "+000", Number: "00000000", Password: "123"})
	newGroup, _ := group.NewGroup(user1, user2)

	if newGroup == nil {
		t.Fatalf("The new group should not be nil.")
	}

	if len(newGroup.Members) < 2 {
		t.Fatalf("The group should have at least 2 users. Gruop.Members: %v", newGroup.Members)
	}
}
