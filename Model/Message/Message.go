package message

import (
	"MessengerService/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	From       user.User
	Content    string
	ReadBy     map[string]primitive.DateTime
	SendedDate primitive.DateTime
}
