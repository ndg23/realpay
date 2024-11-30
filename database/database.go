package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

type Transaction struct {
	ID        string
	Amount    float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDatabase() *Database {
	// Database connection parameters
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Open database connection
	db, err := sql.Open("postgres", "postgres://neondb_owner:PV3xk0GNFspc@ep-dawn-bush-a5nggf59.us-east-2.aws.neon.tech/realpay?sslmode=require")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("Connected to database with DSN: %s", dsn)

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create the transactions table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id VARCHAR(255) PRIMARY KEY,
			amount DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create transactions table: %v", err)
	}

	return &Database{db: db}
}

func (d *Database) SaveTransaction(transaction *Transaction) error {
	query := `
		INSERT INTO transactions (id, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := d.db.Exec(query,
		transaction.ID,
		transaction.Amount,
		transaction.Status,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	return err
}

func (d *Database) GetTransactionByID(id string) (*Transaction, error) {
	query := `
		SELECT id, amount, status, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	transaction := &Transaction{}
	err := d.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return transaction, nil
}

func (d *Database) UpdateTransactionStatus(id, status string) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	result, err := d.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction with ID %s not found", id)
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
