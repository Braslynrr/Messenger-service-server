package dbservice

import (
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type for testing
type DbTest struct {
	users    []*user.User
	groups   []*group.Group
	messages []*message.Message
}

// InsertUser calls dbuser.InsertUser to insert a user in the DB
func (dbs *DbTest) InsertUser(user user.User) (ok bool, err error) {
	dbs.users = append(dbs.users, &user)
	return true, nil
}

// GetUser gets a user from the DB
func (dbs *DbTest) GetUser(localUser user.User) (user *user.User, err error) {
	for _, v := range dbs.users {
		if v.IsEqual(&localUser) {
			return v, err
		}
	}
	return nil, errors.New("user not found")
}

// Login Checks if one user is registed
func (dbs DbTest) Login(localUser user.User) (user *user.User, err error) {
	for _, v := range dbs.users {
		fmt.Println(v)
		fmt.Println(localUser)
		if v.Credentials(&localUser) {
			return v, err
		}
	}
	return nil, errors.New("user can not log in")
}

// CheckGroup checks if chat or group exists
func (dbs DbTest) CheckGroup(user *user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	ID = primitive.NewObjectID()
	for _, v := range dbs.groups {
		count := 0
		for _, y := range v.Members {
			if y.IsEqual(user) {
				count++
			}
			for _, x := range to {
				if y.IsEqual(x) {
					count++
				}
			}
		}
		if count == len(to)+1 {
			return ID, nil
		}
	}
	return ID, errors.New("group not found")
}

// CheckGroup creates a new one
func (dbs *DbTest) CreateGroup(user *user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	ID = primitive.NewObjectID()
	to = append(to, user)
	group, err := group.NewGroup(to...)
	dbs.groups = append(dbs.groups, group)
	return ID, err
}

// SaveMessage Saves message in the DB
func (dbs *DbTest) SaveMessage(message *message.Message) (err error) {
	dbs.messages = append(dbs.messages, message)
	return
}

// GetGroup gets a existing group from db
func (dbs DbTest) GetGroup(ID primitive.ObjectID) (group *group.Group, err error) {
	return
}

// GetAllGroups return all groups of an user
func (dbs DbTest) GetAllGroups(user *user.User) (groups []*group.Group, err error) {
	return dbs.groups, err
}

// GetGroupHistory gets the last messages with a maximun of 20 messages using a date as reference from DB
func (dbs DbTest) GetGroupHistory(groupID primitive.ObjectID, time time.Time) (history []*message.Message, err error) {

	for _, v := range dbs.messages {
		if v.GroupID == groupID {
			history = append(history, v)
		}
	}
	return
}

// UpdateMessageReadBy updates message
func (dbs *DbTest) UpdateMessageReadBy(messageID primitive.ObjectID, localUser user.User) (message message.Message, err error) {
	for _, v := range dbs.messages {
		if v.GroupID == messageID && localUser.IsEqual(v.From) {
			v.IsRead = true
			return *v, err
		}
	}
	return
}
