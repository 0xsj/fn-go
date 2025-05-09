// pkg/common/db/mysql.go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/0xsj/fn-go/pkg/common/config"
	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// MySQLConfig holds MySQL connection configuration
type MySQLConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultMySQLConfig returns default MySQL configuration
func DefaultMySQLConfig() MySQLConfig {
	return MySQLConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "",
		Database:        "app",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// LoadMySQLConfigFromEnv loads MySQL configuration from environment
func LoadMySQLConfigFromEnv(provider config.Provider) MySQLConfig {
	return MySQLConfig{
		Host:            provider.Get("DB_HOST"),
		Port:            provider.GetInt("DB_PORT"),
		Username:        provider.Get("DB_USER"),
		Password:        provider.Get("DB_PASSWORD"),
		Database:        provider.Get("DB_NAME"),
		MaxOpenConns:    provider.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    provider.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: time.Duration(provider.GetInt("DB_CONN_MAX_LIFETIME")) * time.Second,
		ConnMaxIdleTime: time.Duration(provider.GetInt("DB_CONN_MAX_IDLE_TIME")) * time.Second,
	}
}

// MySQLDB implements the DB interface for MySQL
type MySQLDB struct {
	db     *sql.DB
	logger log.Logger
}

// sqlRows implements the Rows interface for sql.Rows
type sqlRows struct {
	rows *sql.Rows
}

func (r *sqlRows) Next() bool {
	return r.rows.Next()
}

func (r *sqlRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *sqlRows) Close() error {
	return r.rows.Close()
}

// sqlRow implements the Row interface for sql.Row
type sqlRow struct {
	row *sql.Row
}

func (r *sqlRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// sqlTransaction implements the Transaction interface for MySQL
type sqlTransaction struct {
	tx *sql.Tx
}

func (t *sqlTransaction) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (t *sqlTransaction) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlRows{rows: rows}, nil
}

func (t *sqlTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return &sqlRow{row: t.tx.QueryRowContext(ctx, query, args...)}
}

// NewMySQLDB creates a new MySQL database connection
func NewMySQLDB(logger log.Logger, config MySQLConfig) (DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	logger.With("host", config.Host).
		With("port", config.Port).
		With("database", config.Database).
		Info("Connecting to MySQL database")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.With("error", err.Error()).Error("Failed to create MySQL database connection")
		return nil, errors.NewDatabaseError("Failed to create database connection", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test connection
	if err := db.Ping(); err != nil {
		logger.With("error", err.Error()).Error("Failed to connect to MySQL database")
		return nil, errors.NewDatabaseError("Failed to connect to database", err)
	}

	logger.Info("Successfully connected to MySQL database")

	return &MySQLDB{
		db:     db,
		logger: logger,
	}, nil
}

// Execute executes a query that doesn't return rows
func (d *MySQLDB) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	// Check if there's a transaction in the context
	if tx, ok := TransactionFromContext(ctx); ok {
		return tx.Execute(ctx, query, args...)
	}

	result, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Query executes a query that returns rows
func (d *MySQLDB) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	// Check if there's a transaction in the context
	if tx, ok := TransactionFromContext(ctx); ok {
		return tx.Query(ctx, query, args...)
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlRows{rows: rows}, nil
}

// QueryRow executes a query that returns a single row
func (d *MySQLDB) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	// Check if there's a transaction in the context
	if tx, ok := TransactionFromContext(ctx); ok {
		return tx.QueryRow(ctx, query, args...)
	}

	return &sqlRow{row: d.db.QueryRowContext(ctx, query, args...)}
}

// WithTransaction executes a function within a transaction
func (d *MySQLDB) WithTransaction(ctx context.Context, fn func(context.Context, Transaction) error) error {
	// Check if there's already a transaction in the context
	if tx, ok := TransactionFromContext(ctx); ok {
		return fn(ctx, tx)
	}

	// Start a new transaction
	sqlTx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.NewDatabaseError("Failed to begin transaction", err)
	}

	tx := &sqlTransaction{tx: sqlTx}
	txCtx := ContextWithTransaction(ctx, tx)

	// Handle panics
	defer func() {
		if p := recover(); p != nil {
			sqlTx.Rollback()
			panic(p) // Re-throw the panic
		}
	}()

	// Execute the function
	if err := fn(txCtx, tx); err != nil {
		if rbErr := sqlTx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %v)", rbErr, err)
		}
		return err
	}

	// Commit the transaction
	if err := sqlTx.Commit(); err != nil {
		return errors.NewDatabaseError("Failed to commit transaction", err)
	}

	return nil
}

// Close closes the database connection
func (d *MySQLDB) Close() error {
	return d.db.Close()
}

// Singleton instance
var (
	mysqlDBInstance DB
	mysqlOnce       sync.Once
	mysqlInitErr    error
)

// GetMySQLSingleton returns a singleton MySQL DB instance
func GetMySQLSingleton(logger log.Logger, config MySQLConfig) (DB, error) {
	mysqlOnce.Do(func() {
		mysqlDBInstance, mysqlInitErr = NewMySQLDB(logger, config)
	})
	return mysqlDBInstance, mysqlInitErr
}