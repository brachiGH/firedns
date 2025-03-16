package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Analytics_DB struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func (a *Analytics_DB) Connect() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Interface on your machine.
	// MongoDB URI and database name
	uri := os.Getenv("MONGO_DB_URI")
	const dbName = "FireDNSanalytics"
	const collectionName = "DNSmessages"

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

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
	// // Find a document
	// var result bson.M
	// err = collection.FindOne(context.Background(), bson.M{"name": "John Doe"}).Decode(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Found a document: %+v\n", result)

	// // Update a document
	// updateResult, err := collection.UpdateOne(
	// 	context.Background(),
	// 	bson.M{"name": "John Doe"},
	// 	bson.M{"$set": bson.M{"age": 31}},
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Matched %v document(s) and modified %v document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// // Delete a document
	// deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"name": "John Doe"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)
}

func (a *Analytics_DB) Disconnect() error {
	if err := a.client.Disconnect(context.Background()); err != nil {
		return fmt.Errorf("error disconnecting: %w", err)
	}
	return nil
}

func (a *Analytics_DB) Update(ip bson.M, doc bson.M) (ID interface{}, err error) {
	// Insert a document
	insertOneResult, err := a.collection.UpdateOne(context.Background(), ip, doc)
	if err != nil {
		return nil, fmt.Errorf("error inserting: %w", err)
	}

	return insertOneResult.UpsertedID, nil
}
