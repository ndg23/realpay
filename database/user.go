// database/user.go
package database

import (
	"context"
	"database/sql"
	"payment-server/model"
)

func (db *Database) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := db.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

func (db *Database) CreateUser(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `
		INSERT INTO users (id, username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (db *Database) CreateAccount(ctx context.Context, tx *sql.Tx, account *model.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, balance, currency, status, kind, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.Balance,
		account.Currency,
		account.Status,
		account.Kind,
		account.CreatedAt,
		account.UpdatedAt,
	)
	return err
}
