package messengermanager

import (
	"MessengerService/groupmanager"
	"MessengerService/message"
	"MessengerService/user"
	"MessengerService/usermanager"
	"sync"

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

// SendMessage initialize the process of sending a message
func (ms *messengerManager) SaveMessage(user *user.User, to []*user.User, message *message.Message) (numbers []string, err error) {
	var ID primitive.ObjectID
	ID, err = groupmanager.HasGroup(user, to)
	if err == nil {
		message.GroupID = ID
		err = groupmanager.SaveMessage(message)
		if err == nil {
			for _, user := range append(to, user) {
				numbers = append(numbers, user.Zone+user.Number)
			}
		}

	}

	return
}

// BroadCastToNumbers broadcast a message to a group of numbers
func (ms *messengerManager) BroadCastToNumbers(numbers []string, message *message.Message) {
	for _, number := range numbers {
		ms.userManager.SendMessageTo(number, message)
	}
}
