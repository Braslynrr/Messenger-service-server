package dbservice

import (
	"MessengerService/dbgroup"
	"MessengerService/dbuser"
	"MessengerService/message"
	"MessengerService/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbService struct {
	connectionLink string
	dbclient       *mongo.Client
	dbcontext      context.Context
	dbdisconnect   context.CancelFunc
}

const mongoLink = "mongodb+srv://Brazza:%v@messengercluster.pgyp5zg.mongodb.net/?retryWrites=true&w=majority"
const passwordFile = "../../DataBase/DBService/mongoPassword.json"

// singleton instance
var (
	instance *dbService
)

var lock = &sync.Mutex{}

// NewDBService creates a unique dbservice instance
func NewDBService() (*dbService, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {

		content, err := os.ReadFile(passwordFile)
		if err != nil {
			return nil, err
		}
		var password map[string]string
		err = json.Unmarshal(content, &password)

		if err != nil {
			return nil, err
		}
		instance = &dbService{connectionLink: fmt.Sprintf(mongoLink, password["password"])}
	}

	return instance, nil
}

// connect connects to the DB
func (dbs *dbService) connect() (err error) {

	dbs.dbcontext, dbs.dbdisconnect = context.WithTimeout(context.Background(), 30*time.Second)
	dbs.dbclient, err = mongo.Connect(dbs.dbcontext, options.Client().ApplyURI(dbs.connectionLink))
	return
}

// close disconnects to the DB
func (dbs *dbService) close() error {

	defer dbs.dbdisconnect()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := dbs.dbclient.Disconnect(dbs.dbcontext); err != nil {
			panic(err)
		}
	}()
	return nil
}

// InsertUser calls dbuser.InsertUser to insert a user in the DB
func (dbs dbService) InsertUser(user user.User) (ok bool, err error) {
	err = dbs.connect()
	user2, err2 := dbuser.GetUser(user, dbs.dbclient, dbs.dbcontext)

	if err2 == nil && user.IsEqual(user2) {
		ok, err = false, errors.New("the user is already registered")
	} else if err == nil {
		ok, err = dbuser.InsertUser(user, dbs.dbclient, dbs.dbcontext)
		dbs.close()
	}
	return
}

// GetUser gets a user from the DB
func (dbs dbService) GetUser(localUser user.User) (user *user.User, err error) {
	err = dbs.connect()

	if err == nil {
		user, err = dbuser.GetUser(localUser, dbs.dbclient, dbs.dbcontext)
		dbs.close()
	}
	return
}

// Login Checks if one user is registed
func (dbs dbService) Login(localUser user.User) (user *user.User, err error) {
	err = dbs.connect()

	if err == nil {
		user, err = dbuser.Login(localUser, dbs.dbclient, dbs.dbcontext)
		dbs.close()
	}
	return
}

// CheckGroup checks if chat or group exists if not, it creates a new one
func (dbs dbService) CheckGroup(user *user.User, to []*user.User) (ID primitive.ObjectID, err error) {
	var groupID any
	var ok bool
	user.UserName = ""
	user.Password = ""
	user.State = ""
	err = dbs.connect()
	if err == nil {
		groupID, err = dbgroup.CheckGroup(user, to, dbs.dbclient, dbs.dbcontext)
		if err != nil {
			groupID, err = dbgroup.CreateGroup(user, to, dbs.dbclient, dbs.dbcontext)
		}
		ID, ok = groupID.(primitive.ObjectID)
		if !ok {
			err = errors.New("it has occured a problem parsing group ID")
		}
	}
	dbs.close()
	return
}

// SaveMessage Saves message in the DB
func (dbs dbService) SaveMessage(message *message.Message) (err error) {
	err = dbs.connect()
	if err == nil {
		err = dbgroup.SaveMessage(message, dbs.dbclient, dbs.dbcontext)
		if err == nil {
			err = dbs.close()
		}
	}
	return
}