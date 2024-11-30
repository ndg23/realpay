package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"payment-server/controllers"
	"payment-server/database"
	"payment-server/middleware"
	"syscall"
	"time"
)

const (
	defaultPort     = "8082"
	shutdownTimeout = 15 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 30 * time.Second
	idleTimeout     = 120 * time.Second
)

func main() {
	// Initialize logger
	log := initLogger()

	// Get port from environment variable or use default
	port := getEnvOrDefault("PORT", defaultPort)

	// Initialize database
	db := database.NewDatabase() // Removed err from assignment since NewDatabase returns 1 value

	// Initialize router and controllers
	router := initRouter(db)

	// Configure server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Str("port", port).Msg("Starting server")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	gracefulShutdown(server, log)
}

func initLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Set log level based on environment
	if os.Getenv("ENV") == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return log
}

func initRouter(db *database.Database) http.Handler {
	router := mux.NewRouter()

	// Initialize controllers
	paymentController := controllers.NewPaymentController(db)
	userService := controllers.NewUserService(db)
	// API versioning middleware
	apiRouter := router.PathPrefix("/v1").Subrouter()

	// Add common middleware
	apiRouter.Use(middleware.RequestLogger)
	apiRouter.Use(middleware.RecoverPanic)
	apiRouter.Use(middleware.ContentTypeJSON)
	//account route
	account := apiRouter.PathPrefix("/account").Subrouter()
	account.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Currency string `json:"currency"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := userService.Register(req.Username, req.Password, req.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	}).Methods(http.MethodPost)

	// Public routes
	payments := apiRouter.PathPrefix("/payments").Subrouter()
	payments.HandleFunc("/init", paymentController.InitializePayment).Methods(http.MethodPost)
	payments.HandleFunc("/{id}/status", paymentController.GetPaymentStatus).Methods(http.MethodGet)

	// Mobile money simulation routes
	payments.HandleFunc("/{id}/confirm", paymentController.ConfirmPayment).Methods(http.MethodPost)
	payments.HandleFunc("/{id}/reject", paymentController.RejectPayment).Methods(http.MethodPost)
	// Removed the cancel payment route since the method doesn't exist

	// Health check endpoint
	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Update this for production
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(router)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func gracefulShutdown(server *http.Server, log zerolog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info().Str("signal", sig.String()).Msg("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Could not gracefully shutdown the server")
	}

	log.Info().Msg("Server stopped")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
