package data

import "time"

type TicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    string `json:"priority"`
}

type TicketResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Link        string    `json:"link"`
	Priority    string    `json:"priority" validate:"required"`
	Created_by  string    `json:"created_by"`
	Created_at  time.Time `json:"created_at"`
}
