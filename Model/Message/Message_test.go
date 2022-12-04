package message_test

import (
	"MessengerService/message"
	"MessengerService/user"
	"testing"
)

// TestNewMessage calls NewMessage to create a new message checking its data
func TestNewMessage(t *testing.T) {
	user, _ := user.NewUser(user.User{Number: "00000000", Zone: "+506", Password: "password"})
	newMessage := message.NewMessage(user, "newMessage")
	if newMessage.From != user {
		t.Fatalf("User and from of message should be the same instance. %v == %v", user, newMessage.From)
	}
	if newMessage.Content != "newMessage" {
		t.Fatalf("Message should be newMessage. Its %v", newMessage.Content)
	}
	if newMessage.From.Password != "" && newMessage.From.State != "" && newMessage.From.UserName != "" {
		t.Fatalf("Password,State and username should be empty. pasword:%v State:%v Username:%v", newMessage.From.Password, newMessage.From.State, newMessage.From.UserName)

	}
}
