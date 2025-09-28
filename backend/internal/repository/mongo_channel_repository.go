package repository

import (
	"context"
	"jobsity-backend/pkg/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoChannelRepository implements ChannelRepository using MongoDB
type MongoChannelRepository struct {
	collection *mongo.Collection
}

// NewMongoChannelRepository creates a new MongoDB channel repository
func NewMongoChannelRepository(collection *mongo.Collection) *MongoChannelRepository {
	return &MongoChannelRepository{
		collection: collection,
	}
}

// Create creates a new channel
func (r *MongoChannelRepository) Create(ctx context.Context, channel *domain.Channel) error {
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, channel)
	if err != nil {
		return err
	}

	// Convert ObjectID to string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		channel.ID = oid.Hex()
	}

	return nil
}

// FindByID finds a channel by ID
func (r *MongoChannelRepository) FindByID(ctx context.Context, id string) (*domain.Channel, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var channel domain.Channel
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// FindByName finds a channel by name
func (r *MongoChannelRepository) FindByName(ctx context.Context, name string) (*domain.Channel, error) {
	var channel domain.Channel
	filter := bson.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// FindAll returns all channels
func (r *MongoChannelRepository) FindAll(ctx context.Context) ([]*domain.Channel, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var channels []*domain.Channel
	for cursor.Next(ctx) {
		var channel domain.Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, err
		}
		channels = append(channels, &channel)
	}

	return channels, nil
}

// Update updates an existing channel
func (r *MongoChannelRepository) Update(ctx context.Context, channel *domain.Channel) error {
	objectID, err := primitive.ObjectIDFromHex(channel.ID)
	if err != nil {
		return err
	}

	channel.UpdatedAt = time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"name":        channel.Name,
		"description": channel.Description,
		"updated_at":  channel.UpdatedAt,
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a channel by ID
func (r *MongoChannelRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
