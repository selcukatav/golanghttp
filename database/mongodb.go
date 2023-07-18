package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_idomitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
}

func addUser(user User) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("GoServer").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)

	if err != nil {
		return err
	}
	return nil

}
func Main() {

	userToAdd := []User{
		{Username: "selcuk", Password: "12345", Role: "admin"},
		{Username: "selcuk1", Password: "123456", Role: "user"},
		{Username: "selcuk2", Password: "123457", Role: "admin"},
	}

	for _, user := range userToAdd {
		err := addUser(user)
		if err != nil {
			log.Printf("an error occured while adding the user")
		} else {
			fmt.Printf("user '%s' added succesfully!!\n", user.Username)
		}

	}

}
