package user

import (
	"errors"
	"regexp"
)

type User struct {
	Number   string `bson:"number"`
	Zone     string `bson:"zone"`
	State    string `bson:"state"`
	UserName string `bson:"username"`
	Password string `bson:"password"`
	socketID string `bson:"-" json:"-"`
}

// NewUser creates a new user, if the parameter is a user which has any golang default property then it will be filled with default values
func NewUser(user User) (*User, error) {
	newUser := &User{Zone: "+000", Number: "00000000", State: "Hi! im using Messeger Service", UserName: "Username"}
	var err error = nil

	if match, err := regexp.MatchString(`\d{8,}`, user.Number); !match {
		if err == nil {
			err = errors.New("Number does not match.")
		}
		return nil, err
	} else {
		newUser.Number = user.Number
	}

	if match, err := regexp.MatchString(`\+\d{3,}`, user.Zone); !match {
		if err == nil {
			err = errors.New("zone does not match.")
		}
		return nil, err
	} else {
		newUser.Zone = user.Zone
	}

	if user.Password == "" {
		err = errors.New("Password cant be empty.")
		return nil, err
	} else {
		newUser.Password = user.Password
	}

	if user.UserName != "" {
		newUser.UserName = user.UserName
	}
	if user.State != "" {
		newUser.State = user.State
	}

	return newUser, nil
}

// IsEqual checks both user are equal by Zone and number
func (user *User) IsEqual(other *User) bool {
	return other != nil && user.Zone == other.Zone && user.Number == other.Number
}

// Credentials check the credentials given and DB infomation are alike.
func (user *User) Credentials(other *User) bool {
	return other != nil && user.Zone == other.Zone && user.Number == other.Number && user.Password == other.Password
}

// SetSocketID sets socket ID
func (user *User) SetSocketID(ID string) {
	user.socketID = ID
}

// GetSocket gets socket ID
func (user *User) GetSocket() string {
	return user.socketID
}
