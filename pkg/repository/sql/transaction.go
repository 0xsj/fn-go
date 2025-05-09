// pkg/repository/sql/transaction.go
package sql

import (
	"context"
	"database/sql"
)

// Context key for transactions
type txContextKey int

const txKey txContextKey = iota

// GetTxFromContext retrieves a transaction from the context if it exists
func GetTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

// WithTransaction executes a function within a transaction
func WithTransaction(ctx context.Context, db *sql.DB, fn func(context.Context) error) error {
	// Check if we're already in a transaction
	if _, ok := GetTxFromContext(ctx); ok {
		// Already in a transaction, just execute the function
		return fn(ctx)
	}
	
	// Start a new transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	// Create a context with the transaction
	txCtx := context.WithValue(ctx, txKey, tx)
	
	// Execute the function
	err = fn(txCtx)
	
	// Handle the result
	if err != nil {
		// Rollback on error
		tx.Rollback()
		return err
	}
	
	// Commit on success
	return tx.Commit()
}

// ExecContext executes a SQL statement within a transaction if one exists
func ExecContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx.ExecContext(ctx, query, args...)
	}
	return db.ExecContext(ctx, query, args...)
}

// QueryContext executes a query within a transaction if one exists
func QueryContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx.QueryContext(ctx, query, args...)
	}
	return db.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query within a transaction if one exists
func QueryRowContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row {
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return db.QueryRowContext(ctx, query, args...)
}