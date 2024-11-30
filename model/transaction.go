package model

import (
	"time"
)

type TransactionRequest struct {
	Amount       float64      `json:"amount"`
	Currency     string       `json:"currency"`
	PayerPhone   string       `json:"payer_phone"`
	Description  string       `json:"description,omitempty"`
	Reference    string       `json:"reference,omitempty"`
	WebhookURL   string       `json:"webhook_url,omitempty"`
	CallbackURLs CallbackURLs `json:"callback_urls,omitempty"`
}

type CallbackURLs struct {
	Success string `json:"success"`
	Error   string `json:"error"`
}

type Transaction struct {
	ID              string    `json:"id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	PayerPhone      string    `json:"payer_phone"`
	Description     string    `json:"description,omitempty"`
	Reference       string    `json:"reference,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	WebhookURL      string    `json:"webhook_url,omitempty"`
	CallbackSuccess string    `json:"callback_success,omitempty"`
	CallbackError   string    `json:"callback_error,omitempty"`
}

type TransactionResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Transaction Transaction `json:"transaction"`
	PaymentURL  string      `json:"payment_url,omitempty"`
}

const (
	StatusPending   = "pending"
	StatusSuccess   = "success"
	StatusError     = "error"
	StatusCancelled = "cancelled"
)