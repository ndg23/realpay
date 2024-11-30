package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"payment-server/database"
	"payment-server/model"
	"payment-server/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type PaymentController struct {
	db *database.Database
}

func NewPaymentController(db *database.Database) *PaymentController {
	return &PaymentController{db: db}
}

func (pc *PaymentController) InitializePayment(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format")
		utils.SendError(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate transaction request
	if err := validateTransactionRequest(&req); err != nil {
		log.Error().Err(err).Msg("Transaction validation failed")
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create transaction
	transaction := createTransactionFromRequest(&req)

	// Save transaction to database
	if err := pc.db.SaveTransaction(&database.Transaction{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.CreatedAt,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to save transaction")
		utils.SendError(w, "Failed to process transaction", http.StatusInternalServerError)
		return
	}

	// Generate payment URL
	paymentURL := generatePaymentURL(transaction.ID)

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Transaction initiated successfully",
		Transaction: transaction,
		PaymentURL:  paymentURL,
	}

	utils.SendSuccess(w, response, http.StatusCreated)
}

func (pc *PaymentController) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	dbTransaction, err := pc.db.GetTransactionByID(transactionID)
	if err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("Transaction retrieval failed")
		utils.SendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if dbTransaction == nil {
		utils.SendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	transaction := model.Transaction{
		ID:        dbTransaction.ID,
		Amount:    dbTransaction.Amount,
		Status:    dbTransaction.Status,
		CreatedAt: dbTransaction.CreatedAt,
		UpdatedAt: dbTransaction.UpdatedAt,
	}

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Status retrieved successfully",
		Transaction: transaction,
	}

	utils.SendSuccess(w, response, http.StatusOK)
}

func (pc *PaymentController) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	dbTransaction, err := pc.db.GetTransactionByID(transactionID)
	if err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("Transaction retrieval failed")
		utils.SendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if dbTransaction == nil {
		utils.SendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	transaction := model.Transaction{
		ID:        dbTransaction.ID,
		Amount:    dbTransaction.Amount,
		Status:    dbTransaction.Status,
		CreatedAt: dbTransaction.CreatedAt,
		UpdatedAt: dbTransaction.UpdatedAt,
	}

	if err := validateTransactionConfirmation(&transaction); err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("Transaction confirmation validation failed")
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update transaction status in database
	if err := pc.db.UpdateTransactionStatus(transactionID, model.StatusSuccess); err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("Failed to update transaction status")
		utils.SendError(w, "Failed to confirm transaction", http.StatusInternalServerError)
		return
	}

	transaction.Status = model.StatusSuccess
	transaction.UpdatedAt = time.Now()

	if transaction.WebhookURL != "" {
		go sendWebhook(r.Context(), transaction)
	}

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Payment confirmed successfully",
		Transaction: transaction,
	}

	utils.SendSuccess(w, response, http.StatusOK)
}

func (pc *PaymentController) RejectPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	dbTransaction, err := pc.db.GetTransactionByID(transactionID)
	if err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("Transaction retrieval failed")
		utils.SendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if dbTransaction == nil {
		utils.SendError(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if err := pc.db.UpdateTransactionStatus(transactionID, model.StatusError); err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("Failed to update transaction status")
		utils.SendError(w, "Failed to reject transaction", http.StatusInternalServerError)
		return
	}

	transaction := model.Transaction{
		ID:        dbTransaction.ID,
		Amount:    dbTransaction.Amount,
		Status:    model.StatusError,
		CreatedAt: dbTransaction.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if transaction.WebhookURL != "" {
		go sendWebhook(r.Context(), transaction)
	}

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Payment rejected successfully",
		Transaction: transaction,
	}

	utils.SendSuccess(w, response, http.StatusOK)
}

func validateTransactionRequest(req *model.TransactionRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("invalid amount: must be positive")
	}
	if req.PayerPhone == "" {
		return fmt.Errorf("payer phone number is required")
	}
	if len(req.PayerPhone) < 10 {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

func createTransactionFromRequest(req *model.TransactionRequest) model.Transaction {
	now := time.Now()
	return model.Transaction{
		ID:              utils.GenerateTransactionID(),
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          model.StatusPending,
		PayerPhone:      req.PayerPhone,
		Description:     req.Description,
		Reference:       req.Reference,
		CreatedAt:       now,
		UpdatedAt:       now,
		ExpiresAt:       now.Add(15 * time.Minute),
		WebhookURL:      req.WebhookURL,
		CallbackSuccess: req.CallbackURLs.Success,
		CallbackError:   req.CallbackURLs.Error,
	}
}

func generatePaymentURL(transactionID string) string {
	return fmt.Sprintf("http://localhost:8082/pay/%s", transactionID)
}

func validateTransactionConfirmation(transaction *model.Transaction) error {
	if transaction.Status != model.StatusPending {
		return fmt.Errorf("transaction cannot be confirmed in current state: %s", transaction.Status)
	}
	if time.Now().After(transaction.ExpiresAt) {
		return fmt.Errorf("transaction has expired")
	}
	return nil
}

func sendWebhook(ctx context.Context, transaction model.Transaction) {
	payload := map[string]interface{}{
		"event_type": "payment.updated",
		"data":       transaction,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Str("transactionID", transaction.ID).Msg("Failed to marshal webhook payload")
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", transaction.WebhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Error().Err(err).Str("transactionID", transaction.ID).Msg("Failed to create webhook request")
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("transactionID", transaction.ID).Msg("Failed to send webhook")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Error().
			Str("transactionID", transaction.ID).
			Int("statusCode", resp.StatusCode).
			Msg("Webhook request failed")
	}
}