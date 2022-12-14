package groupmanager

import (
	"MessengerService/dbservice"
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendMessage sends a message to the DB
func SaveMessage(message *message.Message) (err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		err = dbs.SaveMessage(message)
	}
	return
}

// HasGroup checks if gruup exists
func CheckGroup(user user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		ID, err = dbs.CheckGroup(&user, to)
	}
	return
}

// CreateGroup  create a new group
func CreateGroup(user user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		ID, err = dbs.CreateGroup(&user, to)
	}
	return
}

// GetGroup gets a existing group
func GetGroup(ID primitive.ObjectID) (group *group.Group, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		group, err = dbs.GetGroup(ID)
	}
	return
}

// GetAllGroups returns all groups of an user
func GetAllGroups(user *user.User) (groups []group.Group, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		groups, err = dbs.GetAllGroups(user)
	}
	return
}

// GetGroupHistory gets the last messages with a maximun of 20 messages usincg a date as reference
func GetGroupHistory(groupID primitive.ObjectID, time time.Time) (history []*message.Message, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		history, err = dbs.GetGroupHistory(groupID, time)
	}
	return
}

// UpdateMessageReadBy updates ReadBy field with number and server time
func UpdateMessageReadBy(messageID primitive.ObjectID, user user.User) (message message.Message, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		message, err = dbs.UpdateMessageReadBy(messageID, user)
	}
	return
}
