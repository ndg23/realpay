package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment-server/model"
	"payment-server/utils"
	"time"

	"github.com/gorilla/mux"
)

func InitializePayment(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Format de requête invalide", http.StatusBadRequest)
		return
	}

	// Validation de base
	if req.Amount <= 0 || req.PayerPhone == "" {
		utils.SendError(w, "Montant ou numéro de téléphone invalide", http.StatusBadRequest)
		return
	}

	// Création de la transaction
	transaction := model.Transaction{
		ID:              utils.GenerateTransactionID(),
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          model.StatusPending,
		PayerPhone:      req.PayerPhone,
		Description:     req.Description,
		Reference:       req.Reference,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(15 * time.Minute),
		WebhookURL:      req.WebhookURL,
		CallbackSuccess: req.CallbackURLs.Success,
		CallbackError:   req.CallbackURLs.Error,
	}

	// Ici, vous sauvegarderiez la transaction dans votre base de données

	// Génération de l'URL de paiement (dans un vrai système, ce serait l'URL USSD ou deep link)
	paymentURL := fmt.Sprintf("http://localhost:8080/pay/%s", transaction.ID)

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Transaction initiée avec succès",
		Transaction: transaction,
		PaymentURL:  paymentURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	// 1. Récupérer la transaction depuis la base de données
	// 2. Vérifier que la transaction est dans un état valide
	// 3. Vérifier le solde du compte
	// 4. Effectuer le transfert
	// 5. Mettre à jour le statut

	// Simulons une confirmation réussie
	transaction := model.Transaction{
		ID:     transactionID,
		Status: model.StatusSuccess,
	}

	// Envoyer le webhook de confirmation
	go sendWebhook(transaction)

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Paiement confirmé avec succès",
		Transaction: transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RejectPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	// Mise à jour du statut de la transaction
	transaction := model.Transaction{
		ID:     transactionID,
		Status: model.StatusError,
	}

	// Envoyer le webhook de rejet
	go sendWebhook(transaction)

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Paiement rejeté",
		Transaction: transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	// Ici vous récupéreriez le statut depuis votre base de données
	transaction := model.Transaction{
		ID:     transactionID,
		Status: model.StatusPending,
	}

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Statut récupéré avec succès",
		Transaction: transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, response, http.StatusOK)
}

func CancelPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["id"]

	// Mise à jour du statut de la transaction
	transaction := model.Transaction{
		ID:     transactionID,
		Status: model.StatusError,
	}

	// Envoyer le webhook de rejet
	go sendWebhook(transaction)

	response := model.TransactionResponse{
		Success:     true,
		Message:     "Paiement annulé",
		Transaction: transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, response, http.StatusOK)
}
