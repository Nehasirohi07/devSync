package models

import (
	"time"
)

type Task struct {
	ID            int       `json:"id"`
	ProjectID     int       `json:"project_id"`
	AssignedTo    int       `json:"assigned_to"`
	EmployeeName  string    `json:"employee_name,omitempty"`
	EmployeeEmail string    `json:"employee_email,omitempty"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	Progress      int       `json:"progress"`
	CreatedAt     time.Time `json:"created_at"`
}
