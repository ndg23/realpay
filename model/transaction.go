package model

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	PayerPhone      string    `json:"payer_phone"`
	MerchantPhone   string    `json:"merchant_phone"`
	Description     string    `json:"description"`
	Reference       string    `json:"reference"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	WebhookURL      string    `json:"webhook_url"`
	CallbackSuccess string    `json:"callback_success"`
	CallbackError   string    `json:"callback_error"`
}

type TransactionRequest struct {
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	PayerPhone   string  `json:"payer_phone"`
	Description  string  `json:"description"`
	Reference    string  `json:"reference"`
	WebhookURL   string  `json:"webhook_url"`
	CallbackURLs struct {
		Success string `json:"success"`
		Error   string `json:"error"`
	} `json:"callback_urls"`
}

type TransactionResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	Transaction Transaction `json:"transaction,omitempty"`
	PaymentURL  string      `json:"payment_url,omitempty"`
}

const (
	StatusPending = "pending"
	StatusSuccess = "success"
	StatusError   = "error"
)
