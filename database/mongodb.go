package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//bir fonksiyon oluşturup içinde struct tanımla ve client değerini orada tanımla

type Mongo struct {
	Client *mongo.Client
}

func New() (*Mongo, error) {
	var err error

	//connection of the mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("MongoDb connection error", err)
		return nil, err
	}
	//testing if it is working or not
	if err := client.Ping(context.TODO(), nil); err != nil {
		fmt.Println("Mongodb Ping Error", err)
		return nil, err
	}

	fmt.Println("Successfully connected to the MongoDb.")
	return &Mongo{
		Client: client,
	}, nil
}
