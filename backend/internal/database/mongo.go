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

// ConnectMongo connects to MongoDB and returns client
func ConnectMongo(config Config) (*mongo.Client, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully!")
	return client, nil
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
