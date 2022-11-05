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
}

//NewUser creates a new user, if the parameter is a user which has any golang default property then it will be filled with default values
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
