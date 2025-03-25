package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Analytics_DB struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func (a *Analytics_DB) Connect() error {
	// Interface on your machine.
	// MongoDB URI and database name
	uri := os.Getenv("MONGO_DB_URI")
	const dbName = "FireDNSanalytics"
	const collectionName = "DNSmessages"

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	// Connect to MongoDB
	a.client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting: %w", err)
	}

	// Check the connection
	err = a.client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error connecting: connection check failed: %w", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for the collection
	a.collection = a.client.Database(dbName).Collection(collectionName)

	return nil
}

func (a *Analytics_DB) Disconnect() error {
	if err := a.client.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("error disconnecting: %w", err)
	}
	return nil
}

func (a *Analytics_DB) Update(ip bson.M, doc bson.M) (ID interface{}, err error) {
	updateOptions := options.Update().SetUpsert(true)
	insertOneResult, err := a.collection.UpdateOne(context.Background(), ip, doc, updateOptions)
	if err != nil {
		return nil, fmt.Errorf("error updating: %w", err)
	}

	return insertOneResult.UpsertedID, nil
}
