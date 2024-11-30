package model

import "time"

type Account struct {
	ID       string `json:"id"`
	Balance  int    `json:"balance"`
	Status   string `json:"status"`
	Kind     string `json:"kind"`
	UserID   string `json:"user_id"`
	Currency string `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


