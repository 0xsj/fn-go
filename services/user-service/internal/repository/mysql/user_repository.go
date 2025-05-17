// services/user-service/internal/repository/mysql/user_repository.go
package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/repository"
)

// UserRepository implements repository.UserRepository using MySQL
type UserRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewUserRepository(db *sql.DB, logger log.Logger) repository.UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger.WithLayer("mysql-user-repository"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password, first_name, last_name, phone, 
			role, status, email_verified, preferences, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	preferencesJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return errors.NewInternalError("failed to marshal preferences", err)
	}

	_, err = r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Phone,
		string(user.Role),
		string(user.Status),
		user.EmailVerified,
		preferencesJSON,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "users.username") {
				return errors.NewConflictError("username already exists", err)
			}
			if strings.Contains(err.Error(), "users.email") {
				return errors.NewConflictError("email already exists", err)
			}
			return errors.NewConflictError("user already exists", err)
		}
		return errors.NewDatabaseError("failed to create user", err)
	}
	return nil
}