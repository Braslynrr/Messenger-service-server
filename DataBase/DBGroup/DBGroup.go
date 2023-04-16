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
func CheckGroup(owner *user.User, to []*user.User, client *mongo.Client, ctx context.Context) (id any, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	members := append([]*user.User{owner}, to...)
	filters := bson.M{"members": bson.M{"$size": len(members), "$all": members}}
	result := collection.FindOne(ctx, filters)
	group := &group.Group{}
	err = result.Decode(group)
	return group.ID, err
}

// CreateGroup creates a new chat/group
func CreateGroup(owner *user.User, to []*user.User, client *mongo.Client, ctx context.Context) (id any, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	members := append([]*user.User{owner}, to...)
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
	filters := bson.M{"_id": ID}
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
		member.State = ""
		users = append(users, member)
	}
	return
}

// GetAllGroups gets all groups and its members
func GetAllGroups(localuser *user.User, client *mongo.Client, ctx context.Context) (groups []*group.Group, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	cursor, err := collection.Find(ctx, bson.M{"members": bson.M{"$all": bson.A{localuser}}})
	for cursor.Next(ctx) {
		memberGroup := &group.Group{}
		cursor.Decode(memberGroup)
		memberGroup.Members, err = GetUsersFromGroup(memberGroup.Members, client, ctx)
		groups = append(groups, memberGroup)
	}

	return

}

// GetGroupHistory gets the last messages with a maximun of 20 messages using a date as reference
func GetGroupHistory(groupID primitive.ObjectID, time time.Time, client *mongo.Client, ctx context.Context) (history []*message.Message, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	cursor, err := collection.Find(ctx, bson.M{"groupid": groupID, "sentdate": bson.M{"$lte": time}}, &options.FindOptions{Sort: bson.M{"sentdate": -1}})
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

// GetGroupHistory gets the last messages with a maximun of 20 messages using a date as reference
func UpdateMessageReadBy(messageID primitive.ObjectID, user user.User, client *mongo.Client, ctx context.Context) (err error) {
	user.LeaveMinimalInformation()
	collection := client.Database("Messenger").Collection("Messages")
	_, err = collection.UpdateByID(ctx, messageID, bson.M{"$set": bson.M{"readby": bson.M{user.Zone + user.Number: time.Now()}}})
	return
}

// GetMessage gets a message from DB using an ID
func GetMessage(messageID primitive.ObjectID, client *mongo.Client, ctx context.Context) (message message.Message, err error) {
	collection := client.Database("Messenger").Collection("Messages")
	result := collection.FindOne(ctx, bson.M{"_id": messageID})
	err = result.Decode(&message)
	return
}
