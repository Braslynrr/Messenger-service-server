package group_test

import (
	"MessengerService/group"
	"MessengerService/user"
	"testing"
)

// TestNewGroupFail calls NewGroup checking an error
func TestNewGroupFail(t *testing.T) {
	_, err := group.NewGroup(user.NewUser(user.User{}))
	if err == nil {
		t.Fatalf("TestNewGroup Should return a error. Err: %v", err)
	}
}

// TestNewGroup calls NewGroup checking a new group is correctly created
func TestNewGroup(t *testing.T) {
	newGroup, err := group.NewGroup(user.NewUser(user.User{}), user.NewUser(user.User{}))
	if err != nil {
		t.Fatalf("An error should not be returned. Error: %v", err)
	}
	if newGroup == nil {
		t.Fatalf("The new group should not be nil.")
	}
}

// TestMinimunUsers calls NewGroup checking the new group have at least2 members
func TestMinimunUsers(t *testing.T) {
	newGroup, _ := group.NewGroup(user.NewUser(user.User{}), user.NewUser(user.User{}))

	if newGroup == nil {
		t.Fatalf("The new group should not be nil.")
	}

	if len(newGroup.Members) < 2 {
		t.Fatalf("The group should have at least 2 users. Gruop.Members: %v", newGroup.Members)
	}
}
