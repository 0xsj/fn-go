package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBType string

const (
    MySQL DBType = "mysql"
    PostgreSQL DBType = "postgres"
    MongoDB DBType = "mongodb"
)

type SQLConfig struct {
    Type         DBType
    Host         string
    Port         int
    User         string
    Password     string
    DatabaseName string
    Params       string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}

type NoSQLConfig struct {
    Type         DBType
    URI          string
    DatabaseName string
}

func NewSQLConnection(cfg SQLConfig) (*sqlx.DB, error) {
    var dsn string
    var driver string
    
    switch cfg.Type {
    case MySQL:
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
            cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName, cfg.Params)
        driver = "mysql"
    case PostgreSQL:
        dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
            cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName)
        driver = "postgres"
    default:
        return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
    }
    
    db, err := sqlx.Connect(driver, dsn)
    if err != nil {
        return nil, err
    }
    
    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.MaxLifetime)
    
    return db, nil
}

func NewNoSQLConnection(cfg NoSQLConfig) (*mongo.Database, error) {
    switch cfg.Type {
    case MongoDB:
        clientOptions := options.Client().ApplyURI(cfg.URI)
        
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        client, err := mongo.Connect(ctx, clientOptions)
        if err != nil {
            return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
        }
        
        err = client.Ping(ctx, readpref.Primary())
        if err != nil {
            return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
        }
        
        return client.Database(cfg.DatabaseName), nil
    default:
        return nil, fmt.Errorf("unsupported NoSQL database type: %s", cfg.Type)
    }
}

type Transaction interface {
    Exec(query string, args ...interface{}) (sql.Result, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    Commit() error
    Rollback() error
}

func RunInTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
    tx, err := db.Beginx()
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
    }()
    
    if err := fn(tx); err != nil {
        _ = tx.Rollback()
        return err
    }
    
    return tx.Commit()
}