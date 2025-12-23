package middleware

import (
	"net/http"
	"strings"

	"go-api-server/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks for a valid JWT token in the Authorization header.
// If the token is valid, it sets the user ID in the context and calls the next handler.
// If the token is invalid or missing, it aborts the request with a 401 Unauthorized status.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header format is "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		// Validate the token
		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID in the context so handlers can use it
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		// Call the next handler
		c.Next()
	}
}
