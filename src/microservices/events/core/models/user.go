package models

import (
	"time"
)

type UserEvent struct {
	UserID    int       `json:"user_id" binding:"required"`
	Username  *string   `json:"username,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Action    string    `json:"action" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}
