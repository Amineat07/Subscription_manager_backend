package data

import "time"

type UserRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Role      string `json:"role"`
}

type UserResponse struct {
	ID         int64     `json:"user_ID"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Updated_at time.Time `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
