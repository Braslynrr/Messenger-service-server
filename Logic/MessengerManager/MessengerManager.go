package messengermanager

import (
	"MessengerService/user"
	"MessengerService/usermanager"
	"sync"

	"github.com/gorilla/websocket"
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

func (ms *messengerManager) Login(user user.User) (token string, err error) {
	ok, err := ms.userManager.Login(user)
	if ok != nil && err == nil {
		token, err = ms.userManager.GenerateToken(ok)
	}
	return
}

func (ms *messengerManager) ConnectUser(token string, conn *websocket.Conn) error {
	_, err := ms.userManager.Connect(token, conn)
	return err
}

func (ms *messengerManager) HasTokenAccess(token string) (user *user.User, err error) {
	user, err = ms.userManager.ProcessToken(token)
	return
}
