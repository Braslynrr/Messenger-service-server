package dbgroup

import (
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckGroup checks if a group/chat exists
func CheckGroup(user *user.User, to []*user.User, client *mongo.Client, ctx context.Context) (id any, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	members := append(to, user)
	filters := bson.M{"Members": members}
	result := collection.FindOne(ctx, filters)
	group := &group.Group{}
	err = result.Decode(group)
	return group.ID, err
}

// CreateGroup creates a new chat/group
func CreateGroup(user *user.User, to []*user.User, client *mongo.Client, ctx context.Context) (id any, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	members := append(to, user)
	group, err := group.NewGroup(members...)
	if err == nil {

		var dbgroup *mongo.InsertOneResult
		dbgroup, err = collection.InsertOne(ctx, group)
		id = dbgroup.InsertedID
	}

	return
}

func SaveMessage(message *message.Message, client *mongo.Client, ctx context.Context) (err error) {
	collection := client.Database("Messenger").Collection("Messages")
	_, err = collection.InsertOne(ctx, message)
	return
}
