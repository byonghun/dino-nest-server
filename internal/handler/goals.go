package handler

import (
	"net/http"
	"time"

	"go-api-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateGoalHandler handles the creation of a new savings goal.
func CreateGoalHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal := &models.Goal{
		ID:            uuid.New().String(),
		UserID:        userID.(string),
		Title:         req.Title,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0,
		Duration:      req.Duration,
		StartDate:     time.Now(),
		CreatedAt:     time.Now(),
		Completed:     false,
	}

	// Calculate EndDate based on Duration
	switch req.Duration {
	case models.Weekly:
		goal.EndDate = goal.StartDate.AddDate(0, 0, 7)
	case models.Monthly:
		goal.EndDate = goal.StartDate.AddDate(0, 1, 0)
	case models.Yearly:
		goal.EndDate = goal.StartDate.AddDate(1, 0, 0)
	}

	if err := DB.CreateGoal(goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// GetGoalsHandler retrieves all goals for the authenticated user.
func GetGoalsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goals, err := DB.GetGoalsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve goals"})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// UpdateGoalProgressHandler updates the current amount of a goal.
func UpdateGoalProgressHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goalID := c.Param("id")
	var req models.UpdateGoalProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := DB.GetGoalByID(goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	if goal.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	goal.CurrentAmount += req.Amount
	if goal.CurrentAmount >= goal.TargetAmount {
		goal.Completed = true
		now := time.Now()
		goal.CompletedAt = &now
	}

	if err := DB.UpdateGoal(goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goal"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoalHandler removes a goal.
func DeleteGoalHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	goalID := c.Param("id")
	goal, err := DB.GetGoalByID(goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	if goal.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := DB.DeleteGoal(goalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}
