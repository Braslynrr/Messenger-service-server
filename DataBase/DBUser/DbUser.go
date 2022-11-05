package dbuser

import (
	"MessengerService/user"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(user user.User, client *mongo.Client, ctx context.Context) (bool, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database("Messenger").Collection("Messenger")

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(ctx, user)
	return result != nil, err

}
