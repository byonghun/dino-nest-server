package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go-api-server/internal/database"
	"go-api-server/internal/models"
	"go-api-server/internal/utils"
)

// DB is the global database instance.
// In production, you'd typically use dependency injection instead of a global variable.
// We'll initialize this in main.go and use it across all handlers.
var DB *database.InMemoryDB

// SignupHandler handles user registration requests.
// It creates a new user account with a hashed password and returns a JWT token.
// POST /signup
// Request body: { "email": "user@example.com", "password": "password123" }
// Response: { "token": "jwt-token-here", "user": { "id": "...", "email": "...", "created_at": "..." } }
func SignupHandler(c *gin.Context) {
	// Parse and validate the request body
	var req models.SignupRequest
	
	// ShouldBindJSON automatically validates the request body against our struct tags
	// It checks for required fields, email format, password length, etc.
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return 400 Bad Request if validation fails
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}
	
	// Check if a user with this email already exists
	_, err := DB.GetUserByEmail(req.Email)
	if err == nil {
		// If err is nil, it means we found a user (GetUserByEmail succeeded)
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}
	
	// Hash the password before storing it
	// NEVER store plain text passwords! bcrypt is a secure hashing algorithm
	// The cost parameter (10) controls how secure but slow the hashing is
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		// Return 500 Internal Server Error if hashing fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	
	// Create a new user instance
	user := &models.User{
		// Generate a unique ID using UUID (Universally Unique Identifier)
		ID:        uuid.New().String(),
		Email:     req.Email,
		// Store the hashed password, not the plain text one
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Save the user to the database
	if err := DB.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}
	
	// Generate a JWT token for the new user
	// This allows them to be immediately logged in after signup
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	// Return success response with token and user info
	// Note: We don't include the password in the response
	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// LoginHandler handles user login requests.
// It verifies the email and password, then returns a JWT token if successful.
// POST /login
// Request body: { "email": "user@example.com", "password": "password123" }
// Response: { "token": "jwt-token-here", "user": { "id": "...", "email": "...", "created_at": "..." } }
func LoginHandler(c *gin.Context) {
	// Parse and validate the request body
	var req models.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}
	
	// Look up the user by email
	user, err := DB.GetUserByEmail(req.Email)
	if err != nil {
		// User not found - return 401 Unauthorized
		// Note: We use the same error message for "user not found" and "wrong password"
		// This is a security best practice to prevent email enumeration attacks
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	
	// Verify the password by comparing the hash
	// CompareHashAndPassword checks if the plain text password matches the hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Password doesn't match - return 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	
	// Password is correct! Generate a JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	// Return success response with token and user info
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// LogoutHandler handles user logout requests.
// Note: JWT tokens are stateless, so we can't truly "invalidate" them on the server.
// In a real application, you would either:
// 1. Keep a blacklist of revoked tokens in Redis or a database
// 2. Use short-lived access tokens with refresh tokens
// 3. Let the client just delete the token (client-side logout)
// For this example, we'll just return a success message.
// POST /logout
// Headers: Authorization: Bearer <jwt-token>
// Response: { "message": "Successfully logged out" }
func LogoutHandler(c *gin.Context) {
	// Get the token from the Authorization header
	// Format: "Bearer <token>"
	authHeader := c.GetHeader("Authorization")
	
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No authorization header provided",
		})
		return
	}
	
	// Extract the token part (remove "Bearer " prefix)
	// We expect the header to be in the format: "Bearer <token>"
	var tokenString string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header format",
		})
		return
	}
	
	// Validate the token to make sure it's legitimate
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token: " + err.Error(),
		})
		return
	}
	
	// In a real app, you might:
	// 1. Add the token to a blacklist in Redis with expiration
	// 2. Delete a refresh token from the database
	// 3. Clear server-side session data
	// For now, we'll just return a success message
	// The client should delete the token from their storage (localStorage, cookies, etc.)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
		"user_id": claims.UserID,
	})
}