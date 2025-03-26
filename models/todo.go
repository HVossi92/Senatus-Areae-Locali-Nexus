package models

import (
	"time"
)

// Status represents the current status of a todo item
type Status string

const (
	// Proposed represents a task that has been proposed but not yet approved
	Proposed Status = "PROPOSED"
	// Approved represents a task that has been approved by the senate
	Approved Status = "APPROVED"
	// InProgress represents a task that is currently being worked on
	InProgress Status = "IN_PROGRESS"
	// Completed represents a task that has been completed
	Completed Status = "COMPLETED"
	// Vetoed represents a task that has been vetoed by the consul
	Vetoed Status = "VETOED"
)

// Todo represents a task in the Roman Senate
type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Priority    int       `json:"priority"` // 1 (highest) to 5 (lowest)
	Sponsor     string    `json:"sponsor"`  // The senator who proposed the task
}
