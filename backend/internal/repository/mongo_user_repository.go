package repository

import (
	"context"
	"jobsity-backend/pkg/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserRepository implements UserRepository using MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{
		collection: collection,
	}
}

// FindByEmail finds a user by email
func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *MongoUserRepository) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// Convert ObjectID to string
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}

	return nil
}

// FindByID finds a user by ID
func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *MongoUserRepository) Update(ctx context.Context, user *domain.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"email":    user.Email,
		"password": user.Password,
	}}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a user by ID
func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
