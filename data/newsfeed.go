package data

import "time"

type NewsFeedRequest struct {
	Title       string    `json:"title" validate:"required"`
	Content     string    `json:"content" validate:"required"`
	ImageUrl    string    `json:"image_url"`
	IsPublished bool      `json:"is_published"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

type NewsFeedResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	ImageUrl    string     `json:"image_url"`
	IsPublished bool       `json:"is_published"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedBy   string     `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedBy   *string    `json:"deleted_by"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type UpdateNewsFeedRequest struct {
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	ImageUrl    string     `json:"image_url"`
	IsPublished bool       `json:"is_published"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	UpdatedBy   string     `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
