// pkg/common/db/db.go
package db

import (
	"context"
	"io"
	"time"

	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/errors"
)

// DB is the interface for database operations
type DB interface {
	// Execute executes a query without returning rows
	Execute(ctx context.Context, query string, args ...interface{}) (int64, error)
	
	// Query executes a query that returns rows
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	
	// QueryRow executes a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	
	// BeginTx starts a transaction
	BeginTx(ctx context.Context) (Tx, error)
	
	// Ping checks the database connection
	Ping(ctx context.Context) error
	
	// Close closes the database connection
	Close() error
}

// Tx is the interface for database transactions
type Tx interface {
	// Execute executes a query within a transaction
	Execute(ctx context.Context, query string, args ...interface{}) (int64, error)
	
	// Query executes a query within a transaction
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	
	// QueryRow executes a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	
	// Commit commits the transaction
	Commit() error
	
	// Rollback aborts the transaction
	Rollback() error
}

// Row is the interface for a single database row
type Row interface {
	// Scan copies values from the row into the provided destinations
	Scan(dest ...interface{}) error
}

// Rows is the interface for database query results
type Rows interface {
	io.Closer
	
	// Next advances to the next row
	Next() bool
	
	// Scan copies values from the current row
	Scan(dest ...interface{}) error
	
	// Columns returns the column names
	Columns() ([]string, error)

	Err() error
}

// DatabaseConfig provides common database configuration settings
type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Timeout         time.Duration
}

// DefaultDatabaseConfig returns default database configuration
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "",
		Database:        "app",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		Timeout:         10 * time.Second,
	}
}

// LoadDatabaseConfigFromEnv loads database configuration from environment
func LoadDatabaseConfigFromEnv(provider config.Provider, prefix string) DatabaseConfig {
	if prefix != "" && prefix[len(prefix)-1] != '_' {
		prefix = prefix + "_"
	}
	
	return DatabaseConfig{
		Host:            provider.Get(prefix + "DB_HOST"),
		Port:            provider.GetInt(prefix + "DB_PORT"),
		Username:        provider.Get(prefix + "DB_USER"),
		Password:        provider.Get(prefix + "DB_PASSWORD"),
		Database:        provider.Get(prefix + "DB_NAME"),
		MaxOpenConns:    provider.GetInt(prefix + "DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    provider.GetInt(prefix + "DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: time.Duration(provider.GetInt(prefix+"DB_CONN_MAX_LIFETIME")) * time.Second,
		ConnMaxIdleTime: time.Duration(provider.GetInt(prefix+"DB_CONN_MAX_IDLE_TIME")) * time.Second,
		Timeout:         time.Duration(provider.GetInt(prefix+"DB_TIMEOUT")) * time.Second,
	}
}

// contextKey is used for context values
type contextKey string

// TxContextKey is the key for transaction context values
const TxContextKey contextKey = "transaction"

// GetTxFromContext retrieves a transaction from context
func GetTxFromContext(ctx context.Context) (Tx, bool) {
	tx, ok := ctx.Value(TxContextKey).(Tx)
	return tx, ok
}

// WithTx adds a transaction to a context
func WithTx(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, TxContextKey, tx)
}

// WithTransaction executes a function within a transaction
func WithTransaction(ctx context.Context, db DB, fn func(context.Context) error) error {
	// Check if we already have a transaction
	if _, ok := GetTxFromContext(ctx); ok {
		// Already in a transaction
		return fn(ctx)
	}
	
	// Start a new transaction
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return errors.NewDatabaseError("failed to begin transaction", err)
	}
	
	// Create a context with the transaction
	txCtx := WithTx(ctx, tx)
	
	// Handle panics
	defer func() {
		if p := recover(); p != nil {
			// A panic occurred, rollback
			tx.Rollback()
			panic(p) // Re-panic
		}
	}()
	
	// Execute the function
	err = fn(txCtx)
	
	// Handle the result
	if err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.NewDatabaseError("rollback failed", err).WithField("rollback_error", rbErr.Error())
		}
		return err
	}
	
	// Commit on success
	if err := tx.Commit(); err != nil {
		return errors.NewDatabaseError("commit failed", err)
	}
	
	return nil
}