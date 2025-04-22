// pkg/common/database/mysql.go

package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLConfig struct {
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

func NewMySQLConnection(cfg MySQLConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName, cfg.Params)
	
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
	
	return db, nil
}