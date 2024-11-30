package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"payment-server/model"
)

import "payment-server/database"

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Currency string `json:"currency"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// UserService handles user-related business logic
type UserService struct {
	db *database.Database
}

func NewUserService(db *database.Database) *UserService {
	return &UserService{db: db}
}

// Register handles user registration
func (s *UserService) Register(username, password, currency string) (*model.User, error) {
	// Validate input
	if len(username) < 3 {
		return nil, fmt.Errorf("username must be at least 3 characters long")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters long")
	}

	// Check if username already exists
	exists, err := s.db.CheckUsernameExists(context.Background(), username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	userID := uuid.New().String()
	accountID := uuid.New().String()

	// Create default account for user
	defaultAccount := model.Account{
		ID:       accountID,
		Balance:  0,
		Status:   "active",
		Kind:     "personal",
		UserID:   userID,
		Currency: currency,
	}

	user := &model.User{
		ID:       userID,
		Username: username,
		Password: string(hashedPassword),
		Accounts: []model.Account{defaultAccount},
	}

	// Save user (in a real app, this would be in a database)
	s.db.CreateUser(context.Background(), nil, user)

	return user, nil
}

// RegisterHandler handles HTTP registration requests
func RegisterHandler(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Register user
		user, err := service.Register(req.Username, req.Password, req.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Prepare response
		resp := RegisterResponse{
			ID:       user.ID,
			Username: user.Username,
			Message:  "Registration successful",
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
