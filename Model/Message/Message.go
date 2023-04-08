package message

import (
	"MessengerService/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:",omitempty"`
	GroupID  primitive.ObjectID `json:"groupID,omitempty"`
	From     *user.User
	Content  string
	ReadBy   map[string]*time.Time `json:",omitempty"`
	IsRead   bool                  `bson:"-"`
	SentDate time.Time
}

// NewMessage creates a new message
func NewMessage(from *user.User, content string) (newMessage *Message) {
	newMessage = &Message{From: from, Content: content, ReadBy: make(map[string]*time.Time), SentDate: time.Now()}
	newMessage.From.Password = ""
	newMessage.From.State = ""
	newMessage.From.UserName = ""
	return
}

func (msg *Message) WillSendtoUser(Senduser *user.User) {
	if msg.From.Number != Senduser.Number || msg.From.Zone != Senduser.Zone {
		var isRead bool = msg.ReadBy[Senduser.Zone+Senduser.Number] != nil
		msg.IsRead = isRead
		msg.ReadBy = make(map[string]*time.Time)
	} else {
		msg.IsRead = true
	}
}
