package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"payment-server/model"
)

func sendWebhook(transaction model.Transaction) {
	if transaction.WebhookURL == "" {
		return
	}

	payload := map[string]interface{}{
		"event_type": "payment.updated",
		"data":       transaction,
	}

	jsonPayload, _ := json.Marshal(payload)

	// Envoi de la notification webhook (à implémenter avec retry et gestion d'erreurs)
	_, err := http.Post(transaction.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Erreur d'envoi du webhook: %v", err)
		// Ici vous pourriez implémenter une queue de retry
	}
}
