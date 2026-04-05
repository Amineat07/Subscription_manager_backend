package data

import "time"

type TicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    string `json:"priority"`
}

type TicketUpdateRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Priority    string    `json:"priority"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TicketReplyRequest struct {
	Message   string    `json:"message"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type TicketResponse struct {
	ID          int64                 `json:"id"`
	UserID      int64                 `json:"user_id"`
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Link        string                `json:"link"`
	Priority    string                `json:"priority" validate:"required"`
	Status      string                `json:"status"`
	CreatedBy   string                `json:"created_by"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedBy   string                `json:"updated_by"`
	UpdatedAt   *time.Time            `json:"updated_at"`
	DeletedBy   *string               `json:"deleted_by"`
	DeletedAt   *time.Time            `json:"deleted_at"`
	Replies     []TicketReplyResponse `json:"replies"`
}

type TicketReplyResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TicketID  int64     `json:"ticket_id"`
	Message   string    `json:"title"`
	CreatedBy string    `json:"created_by"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type TicketStatusRequest struct {
	Status string `json:"status" validate:"required"`
}
