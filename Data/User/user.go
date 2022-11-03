package user

type User struct {
	number   string
	zone     string
	state    string
	userName string
}

//NewUser creates a new user, if the parameter is a user which has any golang default property then it will be filled with default values
func NewUser(user User) *User {
	newUser := &User{zone: "+000", number: "00000000", state: "Hi! im using Messeger Service", userName: "Username"}

	if user.number != "" {
		newUser.number = user.number
	}
	if user.zone != "" {
		newUser.zone = user.zone
	}
	if user.userName != "" {
		newUser.userName = user.userName
	}
	if user.number != "" {
		newUser.state = user.state
	}
	return newUser
}
