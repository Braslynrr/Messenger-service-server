package groupmanager

import (
	"MessengerService/dbservice"
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"

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

func GetAllGroups(user *user.User) (groups []group.Group, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		groups, err = dbs.GetAllGroups(user)
	}
	return
}
