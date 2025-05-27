package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Priority  string    `json:"priority"` // low | medium | high
	DueDate   string    `json:"due_date"` // YYYY-MM-DD
	Tags      string    `json:"tags"`     // comma-separated
}
