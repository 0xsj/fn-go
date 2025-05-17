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

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, phone, 
		       role, status, last_login_at, failed_logins, email_verified, 
		       preferences, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	user := &models.User{}
	var preferencesJSON []byte
	var lastLoginAt sql.NullTime
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.Status,
		&lastLoginAt,
		&user.FailedLogins,
		&user.EmailVerified,
		&preferencesJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("user not found", err)
		}
		return nil, errors.NewDatabaseError("failed to get user", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	if len(preferencesJSON) > 0 {
		if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
			return nil, errors.NewInternalError("failed to unmarshal preferences", err)
		}
	}
	return user, nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `
        SELECT id, username, email, password, first_name, last_name, phone, 
               role, status, last_login_at, failed_logins, email_verified, 
               preferences, created_at, updated_at, deleted_at
        FROM users
        WHERE email = ? AND deleted_at IS NULL
    `
    
    user := &models.User{}
    var preferencesJSON []byte
    var lastLoginAt sql.NullTime
    var deletedAt sql.NullTime
    
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Password,
        &user.FirstName,
        &user.LastName,
        &user.Phone,
        &user.Role,
        &user.Status,
        &lastLoginAt,
        &user.FailedLogins,
        &user.EmailVerified,
        &preferencesJSON,
        &user.CreatedAt,
        &user.UpdatedAt,
        &deletedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.NewNotFoundError("user not found", err)
        }
        return nil, errors.NewDatabaseError("failed to get user by email", err)
    }
    
    // Set nullable fields
    if lastLoginAt.Valid {
        user.LastLoginAt = &lastLoginAt.Time
    }
    if deletedAt.Valid {
        user.DeletedAt = &deletedAt.Time
    }
    
    // Deserialize preferences
    if len(preferencesJSON) > 0 {
        if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
            return nil, errors.NewInternalError("failed to unmarshal preferences", err)
        }
    }
    
    return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
    query := `
        SELECT id, username, email, password, first_name, last_name, phone, 
               role, status, last_login_at, failed_logins, email_verified, 
               preferences, created_at, updated_at, deleted_at
        FROM users
        WHERE username = ? AND deleted_at IS NULL
    `
    
    user := &models.User{}
    var preferencesJSON []byte
    var lastLoginAt sql.NullTime
    var deletedAt sql.NullTime
    
    err := r.db.QueryRowContext(ctx, query, username).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Password,
        &user.FirstName,
        &user.LastName,
        &user.Phone,
        &user.Role,
        &user.Status,
        &lastLoginAt,
        &user.FailedLogins,
        &user.EmailVerified,
        &preferencesJSON,
        &user.CreatedAt,
        &user.UpdatedAt,
        &deletedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.NewNotFoundError("user not found", err)
        }
        return nil, errors.NewDatabaseError("failed to get user by username", err)
    }
    
    // Set nullable fields
    if lastLoginAt.Valid {
        user.LastLoginAt = &lastLoginAt.Time
    }
    if deletedAt.Valid {
        user.DeletedAt = &deletedAt.Time
    }
    
    // Deserialize preferences
    if len(preferencesJSON) > 0 {
        if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
            return nil, errors.NewInternalError("failed to unmarshal preferences", err)
        }
    }
    
    return user, nil
}