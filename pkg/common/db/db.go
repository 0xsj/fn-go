// pkg/common/db/db.go
package db

import (
	"context"
)

// DB is an interface for database operations
type DB interface {
	// Execute executes a query that doesn't return rows
	Execute(ctx context.Context, query string, args ...interface{}) (int64, error)
	
	// Query executes a query that returns rows
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	
	// QueryRow executes a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	
	// WithTransaction executes a function within a transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error
	
	// Close closes the database connection
	Close() error
}

// Rows is an interface for database rows
type Rows interface {
	// Next prepares the next row for reading
	Next() bool
	
	// Scan copies the columns in the current row into the values
	Scan(dest ...interface{}) error
	
	// Close closes the rows iterator
	Close() error
}

// Row is an interface for a single database row
type Row interface {
	// Scan copies the columns in the current row into the values
	Scan(dest ...interface{}) error
}

// Transaction is an interface for a database transaction
type Transaction interface {
	// Execute executes a query that doesn't return rows
	Execute(ctx context.Context, query string, args ...interface{}) (int64, error)
	
	// Query executes a query that returns rows
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	
	// QueryRow executes a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
}

// TransactionKey is a context key for transactions
type transactionKey struct{}

// TransactionFromContext retrieves a transaction from the context
func TransactionFromContext(ctx context.Context) (Transaction, bool) {
	tx, ok := ctx.Value(transactionKey{}).(Transaction)
	return tx, ok
}

// ContextWithTransaction returns a context with a transaction
func ContextWithTransaction(ctx context.Context, tx Transaction) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}