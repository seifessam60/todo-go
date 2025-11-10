package dto

import "time"

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low medium high"`
	Category    string `json:"category"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       string  `json:"title" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description"`
	Priority    string  `json:"priority" binding:"omitempty,oneof=low medium high"`
	Category    *string `json:"category"`
	DueDate     *time.Time  `json:"due_date"`
	Completed 	*bool `json:"completed"`
}