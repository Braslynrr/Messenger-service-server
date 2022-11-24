package usermanager

import (
	"MessengerService/dbservice"
	"MessengerService/user"
	"MessengerService/utils"
	"errors"

	"github.com/gorilla/websocket"
)

type UserManager struct {
	tokenList map[string]*user.User
	UserList  map[string]*user.User
}

// NewUserManger creates a new UserManager
func NewUserManger() *UserManager {
	return &UserManager{tokenList: make(map[string]*user.User), UserList: make(map[string]*user.User)}
}

// InsertUser calls DBservice.InsertUser to insert a user to the DB
func (UM *UserManager) InsertUser(user user.User) (ok bool, err error) {

	ok = false
	dbs, err := dbservice.NewDBService()

	if err == nil {
		ok, err = dbs.InsertUser(user)
	}
	return
}

// Login calls dbs.Login checking a user is registered
func (UM *UserManager) Login(user user.User) (ok *user.User, err error) {

	dbs, err := dbservice.NewDBService()

	if err == nil {
		ok, err = dbs.Login(user)
		if err != nil {
			return
		}
	}
	if ok == nil || !user.Credentials(ok) {
		err = errors.New("The given credentials are incorrect.")
	}

	return
}

// Connect  calls login checks login is ok and adds user to userList
func (UM *UserManager) Connect(token string, conn *websocket.Conn) (ok bool, err error) {

	if UM.tokenList[token] == nil {
		return false, errors.New("Token doesn't exist.")
	}
	return true, nil
}

func (UM *UserManager) GenerateToken(user *user.User) (token string, err error) {
	token, err = utils.GenerateToken()
	if err == nil {
		UM.tokenList[token] = user
	}
	return
}

func (UM *UserManager) ProcessToken(token string) (*user.User, error) {
	if user := UM.tokenList[token]; user != nil {
		UM.UserList[user.Zone+user.Number] = user
		return user, nil
	}
	return nil, errors.New("Invalid Token")
}
