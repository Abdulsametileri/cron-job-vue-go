package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDB() (*mongo.Client, error) {
	client, err := mongo.
		Connect(context.Background(),
			options.Client().
				ApplyURI("mongodb://localhost:27017/testReminder"))

	return client, err
}

func SetupCollection(client *mongo.Client, collectionName string) (*mongo.Collection, error) {
	collection := client.Database("testReminder").Collection(collectionName)
	errDrop := collection.Drop(context.Background())
	return collection, errDrop
}
