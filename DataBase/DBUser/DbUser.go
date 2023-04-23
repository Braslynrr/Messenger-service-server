package dbuser

import (
	"MessengerService/user"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertUser inserts one user to the DB
func InsertUser(user user.User, client *mongo.Client, ctx context.Context) (bool, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database("Messenger").Collection("Messenger")

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(ctx, user)
	return result != nil, err

}

// GetUser get one user on the DB
func GetUser(localuser user.User, client *mongo.Client, ctx context.Context) (finalUser *user.User, err error) {

	collection := client.Database("Messenger").Collection("Messenger")

	filters := [2]bson.M{{"number": bson.M{"$eq": localuser.Number}}, {"zone": bson.M{"$eq": localuser.Zone}}}

	finalUser = &user.User{}

	result := collection.FindOne(ctx, bson.M{"$and": filters})

	err = result.Decode(finalUser)

	return
}

// Login checks user is registered in the DB
func Login(localuser user.User, client *mongo.Client, ctx context.Context) (user *user.User, err error) {
	user, err = GetUser(localuser, client, ctx)
	if user.IsEqual(&localuser) {
		return user, err
	}
	return nil, err
}

// UpdateUser updates an user
func UpdateUser(localuser *user.User, client *mongo.Client, ctx context.Context) (err error) {
	collection := client.Database("Messenger").Collection("Messenger")

	filters := bson.D{{Key: "number", Value: localuser.Number}, {Key: "zone", Value: localuser.Zone}}

	update := [4]bson.M{
		{"$set": bson.M{"state": localuser.State}},
		{"$set": bson.M{"password": localuser.Password}},
		{"$set": bson.M{"username": localuser.UserName}},
		{"$set": bson.M{"url": localuser.Url}},
	}
	_, err = collection.UpdateOne(ctx, filters, update)

	return
}
