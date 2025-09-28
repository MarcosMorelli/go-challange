package service

import (
	"context"
	"errors"
	"log"

	"jobsity-backend/internal/repository"
	"jobsity-backend/pkg/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// Login authenticates a user and returns login response
func (s *UserServiceImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return &domain.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	// Find user by email
	userEntity, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.LoginResponse{
				Success: false,
				Message: "Invalid email or password",
			}, nil
		}
		log.Printf("Database error during login: %v", err)
		return &domain.LoginResponse{
			Success: false,
			Message: "Internal server error",
		}, nil
	}

	// Check password (in a real app, you'd hash the password)
	if userEntity.Password != req.Password {
		return &domain.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		}, nil
	}

	// Successful login
	return &domain.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   "fake-jwt-token-12345",
		User:    userEntity,
	}, nil
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	newUser := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetUserByEmail gets a user by email
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// GetUserByID gets a user by ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
