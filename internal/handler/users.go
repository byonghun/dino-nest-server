package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListUsersHandler(c *gin.Context) {
	  // GET all users from the db
		users := DB.GetAllUsers()

		// respond with the list of users and count
		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"count": len(users),
		})
}

// Search endpoint only accounts for email for now
func GetUserByEmailHandler(c *gin.Context) {
	// Get all possible query parameters
	email := c.Query("email")

	// Validate that at least one search parameter is provided
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one search parameter (email) must be provided",
		})
		return
	}

	// Search by email if provided
	user, err := DB.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Return the found user
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")

	user, err := DB.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}