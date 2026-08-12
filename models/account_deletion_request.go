package models

import (
	"time"
)

type AccountDeletionRequest struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	UserName   string     `json:"user_name"`
	UserEmail  string     `json:"user_email"`
	UserRole   string     `json:"user_role"`
	Reason     string     `json:"reason"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy *int       `json:"reviewed_by,omitempty"`
}
