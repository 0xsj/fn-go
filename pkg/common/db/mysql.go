// pkg/common/db/mysql.go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
)

type MySQLConfig struct {
	DatabaseConfig
	ParseTime bool
	Charset   string
}

func DefaultMySQLConfig() MySQLConfig {
	return MySQLConfig{
		DatabaseConfig: DefaultDatabaseConfig(),
		ParseTime:      true,
		Charset:        "utf8mb4",
	}
}

// MySQLDB implements the DB interface for MySQL
type MySQLDB struct {
	db     *sql.DB
	logger log.Logger
}

// MySQLRow implements the Row interface for MySQL
type MySQLRow struct {
	row *sql.Row
}

// Scan implements the Row.Scan method
func (r *MySQLRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// MySQLRows implements the Rows interface for MySQL
type MySQLRows struct {
	rows *sql.Rows
}

// Close implements the Rows.Close method
func (r *MySQLRows) Close() error {
	return r.rows.Close()
}

// Next implements the Rows.Next method
func (r *MySQLRows) Next() bool {
	return r.rows.Next()
}

// Scan implements the Rows.Scan method
func (r *MySQLRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

// Columns implements the Rows.Columns method
func (r *MySQLRows) Columns() ([]string, error) {
	return r.rows.Columns()
}

// Err implements the Rows.Err method
func (r *MySQLRows) Err() error {
	return r.rows.Err()
}

// MySQLTx implements the Tx interface for MySQL
type MySQLTx struct {
	tx *sql.Tx
}

// Execute implements the Tx.Execute method
func (t *MySQLTx) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Query implements the Tx.Query method
func (t *MySQLTx) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &MySQLRows{rows: rows}, nil
}

// QueryRow implements the Tx.QueryRow method
func (t *MySQLTx) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return &MySQLRow{row: t.tx.QueryRowContext(ctx, query, args...)}
}

// Commit implements the Tx.Commit method
func (t *MySQLTx) Commit() error {
	return t.tx.Commit()
}

// Rollback implements the Tx.Rollback method
func (t *MySQLTx) Rollback() error {
	return t.tx.Rollback()
}

// NewMySQLDB creates a new MySQL database connection
func NewMySQLDB(logger log.Logger, config MySQLConfig) (DB, error) {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.ParseTime,
	)

	logger.With("host", config.Host).
		With("port", config.Port).
		With("database", config.Database).
		With("username", config.Username).
		Info("Connecting to MySQL database")

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.With("error", err.Error()).Error("Failed to create MySQL database connection")
		return nil, errors.NewDatabaseError("failed to open database connection", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Verify the connection
	if err := db.PingContext(ctx); err != nil {
		logger.With("error", err.Error()).Error("Failed to connect to MySQL database")
		return nil, errors.NewDatabaseError("failed to connect to database", err)
	}

	logger.Info("Successfully connected to MySQL database")

	return &MySQLDB{
		db:     db,
		logger: logger,
	}, nil
}

// Execute implements the DB.Execute method
func (d *MySQLDB) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	// Check if we're in a transaction
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx.Execute(ctx, query, args...)
	}

	// Execute directly
	result, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Query implements the DB.Query method
func (d *MySQLDB) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	// Check if we're in a transaction
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx.Query(ctx, query, args...)
	}

	// Query directly
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &MySQLRows{rows: rows}, nil
}

// QueryRow implements the DB.QueryRow method
func (d *MySQLDB) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	// Check if we're in a transaction
	if tx, ok := GetTxFromContext(ctx); ok {
		return tx.QueryRow(ctx, query, args...)
	}

	// Query directly
	return &MySQLRow{row: d.db.QueryRowContext(ctx, query, args...)}
}

// BeginTx implements the DB.BeginTx method
func (d *MySQLDB) BeginTx(ctx context.Context) (Tx, error) {
	// Start a new transaction
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &MySQLTx{tx: tx}, nil
}

// Ping implements the DB.Ping method
func (d *MySQLDB) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close implements the DB.Close method
func (d *MySQLDB) Close() error {
	return d.db.Close()
}

// Singleton management
var (
	mysqlInstance DB
	mysqlOnce     sync.Once
	mysqlInitErr  error
)

// GetMySQLSingleton returns a singleton MySQL database instance
func GetMySQLSingleton(logger log.Logger, config MySQLConfig) (DB, error) {
	mysqlOnce.Do(func() {
		mysqlInstance, mysqlInitErr = NewMySQLDB(logger, config)
	})
	return mysqlInstance, mysqlInitErr
}