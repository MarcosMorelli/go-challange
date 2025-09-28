package repository

import (
	"context"
	"jobsity-backend/pkg/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error
}
