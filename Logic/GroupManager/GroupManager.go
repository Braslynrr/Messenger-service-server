package groupmanager

import (
	"MessengerService/dbservice"
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupManager struct {
	dbservice dbservice.DbInterface
}

func NewGroupManager(DB dbservice.DbInterface) *GroupManager {
	return &GroupManager{dbservice: DB}
}

// SendMessage sends a message to the DB
func (gm *GroupManager) SaveMessage(message *message.Message) (err error) {
	err = gm.dbservice.SaveMessage(message)
	return
}

// HasGroup checks if gruup exists
func (gm *GroupManager) CheckGroup(user user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	ID, err = gm.dbservice.CheckGroup(&user, to)
	return
}

// CreateGroup  create a new group
func (gm *GroupManager) CreateGroup(user user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	ID, err = gm.dbservice.CreateGroup(&user, to)
	return
}

// GetGroup gets a existing group
func (gm *GroupManager) GetGroup(ID primitive.ObjectID) (group *group.Group, err error) {
	group, err = gm.dbservice.GetGroup(ID)
	return
}

// GetAllGroups returns all groups of an user
func (gm *GroupManager) GetAllGroups(user *user.User) (groups []*group.Group, err error) {
	groups, err = gm.dbservice.GetAllGroups(user)
	return
}

// GetGroupHistory gets the last messages with a maximun of 20 messages usincg a date as reference
func (gm *GroupManager) GetGroupHistory(groupID primitive.ObjectID, time time.Time) (history []*message.Message, err error) {
	history, err = gm.dbservice.GetGroupHistory(groupID, time)
	return
}

// UpdateMessageReadBy updates ReadBy field with number and server time
func (gm *GroupManager) UpdateMessageReadBy(messageID primitive.ObjectID, user user.User) (message message.Message, err error) {
	message, err = gm.dbservice.UpdateMessageReadBy(messageID, user)
	return
}
