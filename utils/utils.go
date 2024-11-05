package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment-server/model"
	"time"
)

func GenerateTransactionID() string {
	// À implémenter avec un vrai générateur d'ID unique
	return fmt.Sprintf("TRX_%d", time.Now().Unix())
}
func SendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.TransactionResponse{
		Success: false,
		Message: message,
	})
}
func SendSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
func SendWebhook(transaction model.Transaction) {
	if transaction.WebhookURL == "" {
		return
	}
	// Envoyer le webhook
}
