package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"payment-server/controllers"
)

func main() {
	r := mux.NewRouter()

	// Routes publiques
	r.HandleFunc("/v1/payments/init", controllers.InitializePayment).Methods("POST")
	r.HandleFunc("/v1/payments/{id}/status", controllers.GetPaymentStatus).Methods("GET")

	// Routes pour simuler l'interface mobile money
	r.HandleFunc("/v1/payments/{id}/confirm", controllers.ConfirmPayment).Methods("POST")
	r.HandleFunc("/v1/payments/{id}/reject", controllers.RejectPayment).Methods("POST")
	r.HandleFunc("/v1/payments/{id}/cancel", controllers.CancelPayment).Methods("POST")

	// Route pour les webhooks (notifications)
	// r.HandleFunc("/v1/webhooks/payment", controllers.sendWebhook).Methods("POST")

	log.Printf("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
