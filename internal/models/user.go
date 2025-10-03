package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	TelegramID int64    `json:"telegram_id" db:"telegram_id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Language  string    `json:"language" db:"language"`
	IsPremium bool      `json:"is_premium" db:"is_premium"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserState struct {
	UserID int64  `json:"user_id" db:"user_id"`
	State  string `json:"state" db:"state"`
}
