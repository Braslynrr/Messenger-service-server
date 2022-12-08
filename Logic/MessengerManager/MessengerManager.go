package messengermanager

import (
	"MessengerService/group"
	"MessengerService/groupmanager"
	"MessengerService/message"
	"MessengerService/user"
	"MessengerService/usermanager"
	"sync"
	"time"

	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messengerManager struct {
	userManager *usermanager.UserManager
}

// singleton instance
var (
	instance *messengerManager
)

var lock = &sync.Mutex{}

// NewMessengerManager Creates a unique new instance
func NewMessengerManager() (*messengerManager, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {

		instance = &messengerManager{userManager: usermanager.NewUserManger()}

	}

	return instance, nil
}

// InsertUser calls usermanager.InsertUser to insert a user to the DB
func (ms *messengerManager) InsertUser(user user.User) (ok bool, err error) {
	ok, err = ms.userManager.InsertUser(user)
	return
}

// Login check user credentials to return a new token
func (ms *messengerManager) Login(user user.User) (token string, err error) {
	ok, err := ms.userManager.Login(user)
	if ok != nil && err == nil {
		token, err = ms.userManager.GenerateToken(ok)
	}
	return
}

// HasTokenAccess proccess a token an add user to userlist
func (ms *messengerManager) HasTokenAccess(token string) (user *user.User, err error) {
	user, err = ms.userManager.ProcessToken(token)
	return
}

// CheckGroupc checks if a group already exist
func (ms *messengerManager) CheckGroup(user user.User, to []*user.User) (groupID primitive.ObjectID, err error) {
	groupID, err = groupmanager.CheckGroup(user, to)
	return
}

// CreateGroup create a new group in the DB
func (ms *messengerManager) CreateGroup(user user.User, to []*user.User) (groupID primitive.ObjectID, err error) {
	groupID, err = groupmanager.CreateGroup(user, to)
	return
}

// GetGroup gets a group by its identificator
func (ms *messengerManager) GetGroup(groupID primitive.ObjectID) (group *group.Group, err error) {
	group, err = groupmanager.GetGroup(groupID)
	return
}

// SendMessage initialize the process of sending a message
func (ms *messengerManager) SaveMessage(user *user.User, to []*user.User, message *message.Message) (numbers map[socket.SocketId]bool, err error) {
	err = groupmanager.SaveMessage(message)
	if err == nil {
		var tempNumbers []string
		for _, user := range append(to, user) {
			tempNumbers = append(tempNumbers, user.Zone+user.Number)
		}
		numbers = ms.userManager.MapNumbersToSocketID(tempNumbers)
	}

	return
}

// SendToNumber send a message to a group of numbers
func (ms *messengerManager) SendToNumber(conn *socket.Socket, channel string, numbers map[socket.SocketId]bool, message *message.Message) {
	ms.userManager.SendToNumber(conn, channel, numbers, message)
}

// GetAllGroups gets all groups and its member using an user
func (ms *messengerManager) GetAllGroups(user user.User) (groups []group.Group, err error) {
	groups, err = groupmanager.GetAllGroups(&user)
	return
}

// GetGroupHistory gets the last messages with a maximun of 20 messages using a date as reference
func (ms *messengerManager) GetGroupHistory(groupID primitive.ObjectID, time time.Time) (history []*message.Message, err error) {
	history, err = groupmanager.GetGroupHistory(groupID, time)
	return
}

// MessageWasSeenBy sets a message as senn by user
func (ms *messengerManager) MessageWasSeenBy(messageID primitive.ObjectID, user user.User) (message message.Message, err error) {
	message, err = groupmanager.UpdateMessageReadBy(messageID, user)
	return
}

// MapNumberToSocketID Map a User to SocketID if it's online
func (ms *messengerManager) MapNumberToSocketID(user *user.User) (numbers map[socket.SocketId]bool) {
	numbers = ms.userManager.MapNumbersToSocketID([]string{user.Zone + user.Number})
	return
}
