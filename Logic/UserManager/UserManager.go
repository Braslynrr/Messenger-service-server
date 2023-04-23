package usermanager

import (
	"MessengerService/dbservice"
	"MessengerService/user"
	"MessengerService/utils"
	"errors"
	"time"

	"github.com/zishang520/socket.io/socket"
)

type token struct {
	User   *user.User
	tikcer *time.Ticker
}

type UserManager struct {
	tokenList map[string]*token
	UserList  map[string]*user.User
	dbservice dbservice.DbInterface
}

// NewUserManger creates a new UserManager
func NewUserManger(dbs dbservice.DbInterface) *UserManager {
	return &UserManager{tokenList: make(map[string]*token), UserList: make(map[string]*user.User), dbservice: dbs}
}

// InsertUser calls DBservice.InsertUser to insert a user to the DB
func (UM *UserManager) InsertUser(user user.User) (ok bool, err error) {

	ok, err = UM.dbservice.InsertUser(user)
	return
}

// Login calls dbs.Login checking a user is registered
func (UM *UserManager) Login(user user.User) (ok *user.User, err error) {

	ok, err = UM.dbservice.Login(user)
	if err != nil {
		return
	}
	if ok == nil || !user.Credentials(ok) {
		err = errors.New("the given credentials are incorrect")
	}

	return
}

// GenerateToken Generates a new token and assing it to a user
func (UM *UserManager) FakeGenerateToken(user *user.User, stringtoken string) (err error) {
	UM.tokenList[stringtoken] = &token{User: user, tikcer: time.NewTicker(5 * time.Minute)}
	return
}

// GenerateToken Generates a new token and assing it to a user
func (UM *UserManager) GenerateToken(user *user.User) (stringtoken string, err error) {
	stringtoken, err = utils.GenerateToken()
	if err == nil {
		UM.tokenList[stringtoken] = &token{User: user, tikcer: time.NewTicker(5 * time.Minute)}

		go func() {
			for range UM.tokenList[stringtoken].tikcer.C {
				delete(UM.tokenList, stringtoken)
			}
		}()
	}
	return
}

// ProcessToken Process a token and put it on current userList
func (UM *UserManager) ProcessToken(token string) (*user.User, error) {
	if UM.tokenList[token] != nil {
		if user := UM.tokenList[token].User; user != nil {
			UM.UserList[user.Zone+user.Number] = user
			return user, nil
		}
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

	returnedUser, err = UM.dbservice.GetUser(user)
	returnedUser.Password = ""

	return
}

// HasTokenAccess Checks if user has an active token
func (UM *UserManager) HasTokenAccess(user user.User) string {
	for token, tuser := range UM.tokenList {
		if tuser.User.IsEqual(&user) {
			return token
		}
	}
	return ""
}

// ResetTicker resets user's ticker
func (UM *UserManager) ResetTicker(stringtoken string) (err error) {
	if UM.tokenList[stringtoken] != nil {
		UM.tokenList[stringtoken].tikcer.Reset(5 * time.Minute)
		return
	}
	return errors.New("session has expired")
}

// UpdateUser updates an user
func (UM *UserManager) UpdateUser(user *user.User) (err error) {

	err = UM.dbservice.UpdateUser(user)
	if err == nil {
		go func() {
			for token, us := range UM.tokenList {
				if us.User.IsEqual(user) {
					UM.tokenList[token].User = user
					return
				}
			}
		}()
	}

	return
}
