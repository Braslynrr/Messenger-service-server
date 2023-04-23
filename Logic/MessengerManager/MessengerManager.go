package messengermanager

import (
	"MessengerService/dbservice"
	"MessengerService/group"
	"MessengerService/groupmanager"
	"MessengerService/message"
	"MessengerService/user"
	"MessengerService/usermanager"
	"errors"
	"sync"
	"time"

	"github.com/zishang520/socket.io/socket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messengerManager struct {
	userManager  *usermanager.UserManager
	groupManager *groupmanager.GroupManager
	Keys         map[string]string
}

// singleton instance
var (
	instance *messengerManager
)

var lock = &sync.Mutex{}

// NewMessengerManager Creates a unique new instance
func NewMessengerManager(DB dbservice.DbInterface) (*messengerManager, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &messengerManager{userManager: usermanager.NewUserManger(DB), groupManager: groupmanager.NewGroupManager(DB), Keys: map[string]string{}}
	}

	return instance, nil
}

func (ms *messengerManager) IsInitialize() bool {
	return ms.userManager != nil && ms.groupManager != nil && ms.Keys != nil
}

// InsertUser calls usermanager.InsertUser to insert a user to the DB
func (ms *messengerManager) InsertUser(user user.User) (ok bool, err error) {
	ok, err = ms.userManager.InsertUser(user)
	return
}

// Login check user credentials to return a new token
func (ms *messengerManager) FakeLogin(user user.User, token string) (err error) {
	ok, err := ms.userManager.Login(user)
	if ok != nil && err == nil {
		ms.userManager.FakeGenerateToken(&user, token)
	}
	return
}

// Login check user credentials to return a new token
func (ms *messengerManager) Login(user user.User) (token string, err error) {

	ok, err := ms.userManager.Login(user)
	if ok != nil && err == nil {
		if token = ms.userManager.HasTokenAccess(*ok); token == "" {
			token, err = ms.userManager.GenerateToken(ok)
		}

	}
	return
}

// HasTokenAccess proccess a token an add user to userlist
func (ms *messengerManager) HasTokenAccess(token string) (user *user.User, err error) {
	user, err = ms.userManager.ProcessToken(token)
	return
}

// CheckGroupc checks if a group already exist
func (ms *messengerManager) CheckGroup(user user.User, to []*user.User) (groupID primitive.ObjectID, err error) {
	groupID, err = ms.groupManager.CheckGroup(user, to)
	return
}

// CreateGroup create a new group in the DB
func (ms *messengerManager) CreateGroupByUsers(user user.User, to []*user.User) (groupID primitive.ObjectID, err error) {
	groupID, err = ms.groupManager.CreateGroupByUsers(user, to)
	return
}

// CreateGroup create a new group in the DB
func (ms *messengerManager) CreateGroup(group *group.Group) (ouputGroup *group.Group, err error) {
	if len(group.Members) < 2 {
		return nil, errors.New("to create a new chat must be almost two users")
	}
	ouputGroup, err = ms.groupManager.CreateGroup(group)
	return
}

// GetGroup gets a group by its identificator
func (ms *messengerManager) GetGroup(groupID primitive.ObjectID) (group *group.Group, err error) {
	group, err = ms.groupManager.GetGroup(groupID)
	return
}

// SendMessage initialize the process of sending a message
func (ms *messengerManager) SaveMessage(user *user.User, to []*user.User, message *message.Message) (numbers map[socket.SocketId]bool, err error) {
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		defer wait.Done()
		err = ms.groupManager.SaveMessage(message)
	}()

	go func() {
		defer wait.Done()
		var tempNumbers []string
		for _, user := range append(to, user) {
			tempNumbers = append(tempNumbers, user.Zone+user.Number)
		}
		numbers = ms.userManager.MapNumbersToSocketID(tempNumbers)
	}()

	wait.Wait()

	return
}

// GetAllGroups gets all groups and its member using an user
func (ms *messengerManager) GetAllGroups(user user.User) (groups []*group.Group, err error) {
	groups, err = ms.groupManager.GetAllGroups(&user)
	return
}

// GetGroupHistory gets the last messages with a maximun of 20 messages using a date as reference
func (ms *messengerManager) GetGroupHistory(groupID primitive.ObjectID, time time.Time) (history []*message.Message, err error) {
	history, err = ms.groupManager.GetGroupHistory(groupID, time)
	return
}

// MessageWasSeenBy sets a message as senn by user
func (ms *messengerManager) MessageWasSeenBy(messageID primitive.ObjectID, user user.User) (message message.Message, err error) {
	message, err = ms.groupManager.UpdateMessageReadBy(messageID, user)
	return
}

// MapNumberToSocketID Map a User to SocketID if it's online
func (ms *messengerManager) MapNumberToSocketID(user *user.User) *socket.SocketId {
	numbers := ms.userManager.MapNumbersToSocketID([]string{user.Zone + user.Number})
	for si := range numbers {
		return &si
	}
	return nil
}

// GetUser gets an user from userManager
func (ms *messengerManager) GetUser(user user.User) (returneduser *user.User, err error) {
	returneduser, err = ms.userManager.GetUser(user)
	return
}

// ResetUserTime resets the user ticker time
func (ms *messengerManager) ResetUserTime(token string) error {
	return ms.userManager.ResetTicker(token)
}

// MapUsersToSocketsID maps a users into socketsID list
func (ms *messengerManager) MapUsersToSocketsID(users []*user.User) map[socket.SocketId]bool {
	var numbers []string
	for _, us := range users {
		numbers = append(numbers, us.Zone+us.Number)
	}
	return ms.userManager.MapNumbersToSocketID(numbers)
}

func (ms *messengerManager) UpdateUser(user *user.User) error {
	return ms.userManager.UpdateUser(user)
}
