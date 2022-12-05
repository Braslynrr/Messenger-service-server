package message_test

import (
	"MessengerService/message"
	"MessengerService/user"
	"testing"
	"time"
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

// TestWillSendtoUser test WillSendtoUser checking owner recieves the correct data
func TestWillSendtoUser(t *testing.T) {
	user, _ := user.NewUser(user.User{Number: "00000000", Zone: "+506", Password: "password"})
	newMessage := message.NewMessage(user, "newMessage")
	newMessage.WillSendtoUser(user)
	if newMessage.IsRead != true {
		t.Fatalf("IsRead should be true. NewMessage %v", newMessage)
	}
	var time time.Time
	newMessage.ReadBy[user.Zone+user.Number] = &time
	if len(newMessage.ReadBy) == 0 {
		t.Fatalf("ReadBy should have one reader. NewMessage: %v", newMessage)
	}
}

// TestWillSendtoUserNotOwner test WillSendtoUser checking a not owner recieves the correct data
func TestWillSendtoUserNotOwner(t *testing.T) {
	user2, _ := user.NewUser(user.User{Number: "00000000", Zone: "+507", Password: "password"})
	user, _ := user.NewUser(user.User{Number: "00000000", Zone: "+506", Password: "password"})
	newMessage := message.NewMessage(user, "newMessage")
	newMessage.WillSendtoUser(user2)
	if newMessage.IsRead != false {
		t.Fatalf("IsRead should be false. NewMessage: %v", newMessage)
	}
	var time time.Time
	newMessage.ReadBy[user.Zone+user.Number] = &time
	if len(newMessage.ReadBy) == 0 {
		t.Fatalf("ReadBy should have one reader. NewMessage: %v", newMessage)
	}
}

// TestWillSendtoUserNotOwner test WillSendtoUser checking a not owner that read the message recieves the correct data
func TestWillSendtoUserNotOwnerIsRead(t *testing.T) {
	user2, _ := user.NewUser(user.User{Number: "00000000", Zone: "+507", Password: "password"})
	user, _ := user.NewUser(user.User{Number: "00000000", Zone: "+506", Password: "password"})
	newMessage := message.NewMessage(user, "newMessage")
	var time time.Time
	newMessage.ReadBy[user2.Zone+user2.Number] = &time
	newMessage.WillSendtoUser(user2)
	if newMessage.IsRead != true {
		t.Fatalf("IsRead should be true. NewMessage: %v", newMessage)
	}
	if len(newMessage.ReadBy) != 0 {
		t.Fatalf("ReadBy should not have readers. NewMessage: %v len(NewMessage):%v", newMessage.ReadBy, len(newMessage.ReadBy))
	}
}
