package models

import "time"

// User represents a user in the system.
// This struct defines the structure of user data stored in our in-memory database.
type User struct {
	// ID is the unique identifier for the user
	// In a real database, this would be auto-generated
	ID string `json:"id"`
	
	// Email is the user's email address, used for login
	// This must be unique across all users
	Email string `json:"email"`
	
	// Password stores the hashed password (never store plain text passwords!)
	// We'll use bcrypt to hash passwords before storing
	Password string `json:"-"` // json:"-" means this field won't be included in JSON responses
	
	// CreatedAt tracks when the user account was created
	CreatedAt time.Time `json:"created_at"`
	
	// UpdatedAt tracks when the user account was last modified
	UpdatedAt time.Time `json:"updated_at"`
}

// SignupRequest represents the data required for user registration.
// This is what we expect to receive in the request body for /signup
type SignupRequest struct {
	// Email is the user's email address
	Email string `json:"email" binding:"required,email"`
	
	// Password is the user's chosen password
	// binding:"required,min=6" means this field is required and must be at least 6 characters
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the data required for user login.
// This is what we expect to receive in the request body for /login
type LoginRequest struct {
	// Email is the user's email address
	Email string `json:"email" binding:"required,email"`
	
	// Password is the user's password
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response sent back after successful authentication.
// This is returned to the client after successful signup or login
type AuthResponse struct {
	// Token is the JWT token that the client will use for authenticated requests
	Token string `json:"token"`
	
	// User contains the user's basic information (without sensitive data like password)
	User UserResponse `json:"user"`
}

// UserResponse represents user data that is safe to send to clients.
// Note: We don't include the password field here for security
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
