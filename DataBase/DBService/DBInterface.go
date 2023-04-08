package dbservice

import (
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DbInterface interface {
	InsertUser(user user.User) (ok bool, err error)
	GetUser(localUser user.User) (user *user.User, err error)
	Login(localUser user.User) (user *user.User, err error)
	CheckGroup(user *user.User, to []*user.User) (ID primitive.ObjectID, err error)
	CreateGroup(user *user.User, to []*user.User) (ID primitive.ObjectID, err error)
	SaveMessage(message *message.Message) (err error)
	GetGroup(ID primitive.ObjectID) (group *group.Group, err error)
	GetAllGroups(user *user.User) (groups []*group.Group, err error)
	GetGroupHistory(groupID primitive.ObjectID, time time.Time) (history []*message.Message, err error)
	UpdateMessageReadBy(messageID primitive.ObjectID, localUser user.User) (message message.Message, err error)
}
