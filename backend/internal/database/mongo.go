package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds database configuration
type Config struct {
	URI      string
	Database string
}

// ConnectMongo connects to MongoDB and returns client and collection
func ConnectMongo(config Config) (*mongo.Client, *mongo.Collection, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Get the users collection
	collection := client.Database(config.Database).Collection("users")

	log.Println("Connected to MongoDB successfully!")
	return client, collection, nil
}

// DisconnectMongo disconnects from MongoDB
func DisconnectMongo(client *mongo.Client) {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client.Disconnect(ctx)
		log.Println("Disconnected from MongoDB")
	}
}
