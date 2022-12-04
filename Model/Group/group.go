package group

import (
	"MessengerService/user"
	"errors"
)

type Group struct {
	ID          any `bson:"_id,omitempty"`
	Members     []*user.User
	GroupName   string
	Description string
	IsChat      bool
	Admins      []*user.User
}

// NewGroup creates a new Group, 2 members are needed
func NewGroup(users ...*user.User) (newGroup *Group, err error) {

	newGroup = &Group{Members: make([]*user.User, 0), IsChat: true}
	err = nil
	for _, user := range users {
		if newGroup.Admins == nil {
			newGroup.Admins = append(newGroup.Admins, user)
		}
		newGroup.Members = append(newGroup.Members, user)
	}

	if len(newGroup.Members) < 2 {
		newGroup = nil
		err = errors.New("there aren't enough users")
	} else if len(newGroup.Members) > 3 {
		newGroup.IsChat = false
		for _, member := range newGroup.Members {
			newGroup.GroupName += member.UserName + " "
		}
	}
	return newGroup, err
}
