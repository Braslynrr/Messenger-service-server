package messengermanager

import (
	"MessengerService/user"
	"MessengerService/usermanager"
	"sync"
)

type messengerManager struct {
}

// singleton instance
var (
	instance *messengerManager
)

var lock = &sync.Mutex{}

func NewMessengerManager() (*messengerManager, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {

		instance = &messengerManager{}
	}

	return instance, nil
}

func (ms *messengerManager) InsertUser(user user.User) (ok bool, err error) {
	ok, err = usermanager.InsertUser(user)
	return
}
