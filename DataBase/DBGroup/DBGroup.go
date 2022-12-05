package dbgroup

import (
	"MessengerService/group"
	"MessengerService/message"
	"MessengerService/user"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckGroup checks if a group/chat exists
func CheckGroup(user *user.User, to []*user.User, client *mongo.Client, ctx context.Context) (id any, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	members := append(to, user)
	filters := bson.M{"members": bson.M{"$size": len(members), "$all": members}}
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

// SaveMessage save the message into messages
func SaveMessage(message *message.Message, client *mongo.Client, ctx context.Context) (err error) {
	collection := client.Database("Messenger").Collection("Messages")
	ID, err := collection.InsertOne(ctx, message)
	if err == nil {
		message.ID = ID.InsertedID.(primitive.ObjectID)
	}
	return
}

// GetGroup gets a group and its members
func GetGroup(ID primitive.ObjectID, client *mongo.Client, ctx context.Context) (serverGroup *group.Group, err error) {
	serverGroup = &group.Group{}
	collection := client.Database("Messenger").Collection("Messages")
	filters := bson.M{"_ID": ID}
	result := collection.FindOne(ctx, filters)
	err = result.Decode(serverGroup)
	if err == nil {
		serverGroup.Members, err = GetUsersFromGroup(serverGroup.Members, client, ctx)
	}
	return
}

// GetUsersFromGroup gets all user from an group
func GetUsersFromGroup(members []*user.User, client *mongo.Client, ctx context.Context) (users []*user.User, err error) {
	collection := client.Database("Messenger").Collection("Messenger")
	for _, member := range members {
		filters := bson.M{"zone": member.Zone, "number": member.Number}
		result := collection.FindOne(ctx, filters)
		err = result.Decode(member)
		member.Password = ""
		users = append(users, member)
	}
	return
}

// GetAllGroups gets all groups and its members
func GetAllGroups(localuser *user.User, client *mongo.Client, ctx context.Context) (groups []group.Group, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	cursor, err := collection.Find(ctx, bson.M{"members": bson.M{"$all": bson.A{localuser}}})
	for cursor.Next(ctx) {
		memberGroup := &group.Group{}
		cursor.Decode(memberGroup)
		memberGroup.Members, err = GetUsersFromGroup(memberGroup.Members, client, ctx)
		groups = append(groups, *memberGroup)
	}

	return

}

// GetGroupHistory gets the last messages with a maximun of 20 messages using a date as reference
func GetGroupHistory(groupID primitive.ObjectID, time time.Time, client *mongo.Client, ctx context.Context) (history []*message.Message, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	cursor, err := collection.Find(ctx, bson.M{"groupid": groupID, "sendeddate": bson.M{"$lte": time}}, &options.FindOptions{Sort: bson.M{"sendeddate": -1}})
	if err == nil {
		for cursor.Next(ctx) && len(history) < 20 {
			message := &message.Message{}
			err = cursor.Decode(&message)
			if err == nil {
				history = append(history, message)
			}
		}
	}
	return
}
