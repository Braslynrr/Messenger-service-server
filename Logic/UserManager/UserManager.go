package usermanager

import (
	"MessengerService/dbservice"
	"MessengerService/user"
	"MessengerService/utils"
	"errors"
	"time"

	"github.com/zishang520/socket.io/socket"
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
		err = errors.New("the given credentials are incorrect")
	}

	return
}

// GenerateToken Generates a new token and assing it to a user
func (UM *UserManager) FakeGenerateToken(user *user.User, token string) (err error) {
	UM.tokenList[token] = user
	return
}

// GenerateToken Generates a new token and assing it to a user
func (UM *UserManager) GenerateToken(user *user.User) (token string, err error) {
	token, err = utils.GenerateToken()
	if err == nil {
		UM.tokenList[token] = user
		ticker := time.NewTicker(5 * time.Second)

		go func() {
			for range ticker.C {
				delete(UM.tokenList, token)
			}
		}()
	}
	return
}

// ProcessToken Process a token and put it on current userList
func (UM *UserManager) ProcessToken(token string) (*user.User, error) {
	if user := UM.tokenList[token]; user != nil {
		UM.UserList[user.Zone+user.Number] = user
		return user, nil
	}

	return nil, errors.New("invalid token")
}

// MapNumbersToSocketID Map number to socketsID
func (UM *UserManager) MapNumbersToSocketID(numbers []string) (numberMap map[socket.SocketId]bool) {
	numberMap = make(map[socket.SocketId]bool, 0)
	for _, number := range numbers {
		if UM.UserList[number] != nil {
			numberMap[UM.UserList[number].GetSocket()] = true
		}
	}
	return numberMap
}

// GetUser gets an user from DB
func (UM *UserManager) GetUser(user user.User) (returnedUser *user.User, err error) {
	dbs, err := dbservice.NewDBService()

	if err == nil {
		returnedUser, err = dbs.GetUser(user)
		returnedUser.Password = ""
	}
	return
}
