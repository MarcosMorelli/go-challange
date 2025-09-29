package repository

import (
	"context"
	"jobsity-backend/pkg/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoMessageRepository implements MessageRepository using MongoDB
type MongoMessageRepository struct {
	collection *mongo.Collection
}

// NewMongoMessageRepository creates a new MongoDB message repository
func NewMongoMessageRepository(collection *mongo.Collection) *MongoMessageRepository {
	return &MongoMessageRepository{
		collection: collection,
	}
}

// Create creates a new message
func (r *MongoMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	message.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	// Convert ObjectID to string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		message.ID = oid.Hex()
	}

	return nil
}

// FindByID finds a message by ID
func (r *MongoMessageRepository) FindByID(ctx context.Context, id string) (*domain.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var message domain.Message
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// FindByChannelID finds all messages for a specific channel
func (r *MongoMessageRepository) FindByChannelID(ctx context.Context, channelID string, limit int) ([]*domain.Message, error) {
	filter := bson.M{"channel_id": channelID}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	for cursor.Next(ctx) {
		var message domain.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

// Update updates an existing message
func (r *MongoMessageRepository) Update(ctx context.Context, message *domain.Message) error {
	objectID, err := primitive.ObjectIDFromHex(message.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"content": message.Content,
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a message by ID
func (r *MongoMessageRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
