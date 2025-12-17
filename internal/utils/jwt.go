package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is the secret key used to sign and verify JWT tokens.
// IMPORTANT: In production, this should be:
// 1. Loaded from environment variables (not hardcoded)
// 2. A long, random, secure string
// 3. Kept secret and never committed to version control
// For this learning example, we're hardcoding it, but DON'T do this in real apps!
var JWTSecret = []byte("your-secret-key-change-this-in-production")

// JWTClaims represents the claims (data) stored in our JWT token.
// Claims are the payload of the JWT - the information we want to encode.
type JWTClaims struct {
	// UserID is the unique identifier of the authenticated user
	UserID string `json:"user_id"`
	
	// Email is the user's email address
	Email string `json:"email"`
	
	// RegisteredClaims includes standard JWT fields like expiration time
	// This is provided by the jwt library and includes fields like:
	// - ExpiresAt: when the token expires
	// - IssuedAt: when the token was created
	// - NotBefore: when the token becomes valid
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a user.
// This function is called after successful login or signup.
// Parameters:
//   - userID: the unique identifier of the user
//   - email: the user's email address
// Returns:
//   - string: the signed JWT token as a string
//   - error: nil if successful, error if token generation fails
func GenerateJWT(userID, email string) (string, error) {
	// Set the token expiration time to 24 hours from now
	// In production, you might want shorter expiration times (e.g., 15 minutes)
	// and use refresh tokens for longer sessions
	expirationTime := time.Now().Add(24 * time.Hour)
	
	// Create the claims (payload) for the token
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			// Set when the token expires
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// Set when the token was issued (now)
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// NotBefore: token is valid immediately
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	// Create a new token with the claims and the HS256 signing method
	// HS256 = HMAC with SHA-256, a symmetric signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token with our secret key to produce the final JWT string
	// This creates the signature part of the JWT
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	
	// Return the complete JWT token string
	return tokenString, nil
}

// ValidateJWT verifies a JWT token and extracts the claims from it.
// This function is called on protected routes to authenticate requests.
// Parameters:
//   - tokenString: the JWT token string to validate
// Returns:
//   - *JWTClaims: pointer to the claims if token is valid
//   - error: nil if valid, error describing why validation failed
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	// Initialize claims struct to store the decoded data
	claims := &JWTClaims{}
	
	// Parse the token string and verify its signature
	// The callback function provides the secret key for verification
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is what we expect (HS256)
		// This prevents attacks where someone changes the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return our secret key for signature verification
		return JWTSecret, nil
	})
	
	// Check if there was an error during parsing
	if err != nil {
		return nil, err
	}
	
	// Check if the token is valid (signature verified and not expired)
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	// Return the extracted claims
	return claims, nil
}

// ExtractUserIDFromToken is a convenience function to get just the user ID from a token.
// This is useful when you only need the user ID and don't care about other claims.
// Parameters:
//   - tokenString: the JWT token string
// Returns:
//   - string: the user ID extracted from the token
//   - error: nil if successful, error if token is invalid
func ExtractUserIDFromToken(tokenString string) (string, error) {
	// Validate the token and get the claims
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}
	
	// Return just the user ID from the claims
	return claims.UserID, nil
}
