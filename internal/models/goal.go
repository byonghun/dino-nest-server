package models

import "time"

type GoalDuration string

const (
    Weekly  GoalDuration = "weekly"
    Monthly GoalDuration = "monthly"
    Yearly  GoalDuration = "yearly"
)

type Goal struct {
    ID            string       `json:"id"`
    UserID        string       `json:"user_id"`
    Title         string       `json:"title"`
    TargetAmount  float64      `json:"target_amount"`
    CurrentAmount float64      `json:"current_amount"`
    Duration      GoalDuration `json:"duration"`
    StartDate     time.Time    `json:"start_date"`
    EndDate       time.Time    `json:"end_date"`
    Completed     bool         `json:"completed"`
    CompletedAt   *time.Time   `json:"completed_at,omitempty"`
    CreatedAt     time.Time    `json:"created_at"`
}

type CreateGoalRequest struct {
    Title        string       `json:"title" binding:"required"`
    TargetAmount float64      `json:"target_amount" binding:"required,gt=0"`
    Duration     GoalDuration `json:"duration" binding:"required,oneof=weekly monthly yearly"`
}

type UpdateGoalProgressRequest struct {
    Amount float64 `json:"amount" binding:"required,gt=0"`
}
