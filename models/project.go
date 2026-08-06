package models

import (
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ManagerID   int       `json:"manager_id"`
	CreatedAt   time.Time `json:"created_at"`
}
