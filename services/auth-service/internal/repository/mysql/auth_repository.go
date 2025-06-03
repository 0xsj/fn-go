// services/auth-service/internal/repository/mysql/auth_repository.go
package repository

import (
	"database/sql"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/services/auth-service/internal/repository"
)

// AuthRepository implements the composite AuthRepository interface
type AuthRepository struct {
	repository.TokenRepository
	repository.SessionRepository
	repository.PermissionRepository
	repository.RolePermissionRepository
	
	db     *sql.DB
	logger log.Logger
}

// NewAuthRepository creates a new composite auth repository
func NewAuthRepository(db *sql.DB, logger log.Logger) repository.AuthRepository {
	return &AuthRepository{
		TokenRepository:      NewTokenRepository(db, logger),
		SessionRepository:    NewSessionRepository(db, logger),
		PermissionRepository: NewPermissionRepository(db, logger),
		RolePermissionRepository: NewRolePermissionRepository(db, logger),
		db:     db,
		logger: logger.WithLayer("mysql-auth-repository"),
	}
}

// Transaction method for the composite repository
func (r *AuthRepository) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

// GetDB returns the database connection for advanced operations
func (r *AuthRepository) GetDB() *sql.DB {
	return r.db
}