package message

import (
	"MessengerService/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:",omitempty"`
	GroupID    primitive.ObjectID `json:"groupID,omitempty"`
	From       *user.User
	Content    string
	ReadBy     map[string]time.Time
	SendedDate time.Time
}

// NewMessage creates a new message
func NewMessage(from *user.User, content string) (newMessage *Message) {
	newMessage = &Message{From: from, Content: content, ReadBy: make(map[string]time.Time), SendedDate: time.Now()}
	newMessage.From.Password = ""
	newMessage.From.State = ""
	newMessage.From.UserName = ""
	return
}
