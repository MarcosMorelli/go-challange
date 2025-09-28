package domain

import "time"

// User represents a user in the system
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// LoginRequest represents the login request structure
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response structure
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// CreateUserRequest represents the create user request structure
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
