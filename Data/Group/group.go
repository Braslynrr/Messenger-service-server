package group

import (
	"MessengerService/user"
	"errors"
)

type Group struct {
	Members     []*user.User
	Name        string
	Description string
	IsChat      bool
	Admins      []*user.User
}

// NewGroup creates a new Group, 2 members are needed
func NewGroup(users ...*user.User) (newGroup *Group, err error) {

	newGroup = &Group{Members: make([]*user.User, 0)}
	err = nil

	for _, user := range users {
		if newGroup.Admins == nil {
			newGroup.Admins = append(newGroup.Admins, user)
		}
		newGroup.Members = append(newGroup.Members, user)
	}

	if len(newGroup.Members) < 2 {
		newGroup = nil
		err = errors.New("There aren't enough users")
	}
	return newGroup, err
}
