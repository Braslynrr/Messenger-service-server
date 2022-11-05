package usermanager

import (
	"MessengerService/dbservice"
	"MessengerService/user"
)

func InsertUser(user user.User) (ok bool, err error) {

	ok = false
	dbs, err := dbservice.NewDBService()

	if err == nil {
		ok, err = dbs.InsertUser(user)
	}
	return
}
