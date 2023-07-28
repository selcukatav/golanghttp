package middlewares

import (
	"context"
	"fmt"
	"goserver/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Authenticate(client *database.Mongo, username, password, role string) (*User, error) {

	if client == nil {
		fmt.Println("client değişkeni nil.")

	}

	collection := client.Client.Database("GoServer").Collection("users")

	filter := bson.M{"username": username, "password": password, "role": role}

	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("username or password is worng")
		}
		return nil, fmt.Errorf("error occured in authenticate")
	}
	return &user, nil

}
