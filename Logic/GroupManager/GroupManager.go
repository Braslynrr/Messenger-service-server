package groupmanager

import (
	"MessengerService/dbservice"
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

// HasGroup checks if gruup exists otherwise it will create a new group
func HasGroup(user *user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	dbs, err := dbservice.NewDBService()
	if err == nil {
		ID, err = dbs.CheckGroup(user, to)
	}
	return
}
