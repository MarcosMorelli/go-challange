package service

import (
	"context"
	"jobsity-backend/pkg/domain"
)

// UserService defines the interface for user business logic
type UserService interface {
	// Login authenticates a user and returns login response
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)

	// GetUserByEmail gets a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// GetUserByID gets a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}
