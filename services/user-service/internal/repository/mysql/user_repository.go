// services/user-service/internal/repository/mysql/user_repository.go
package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/0xsj/fn-go/pkg/common/db"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/domain"
	"github.com/0xsj/fn-go/services/user-service/internal/repository"
)

// UserRepository implements repository.UserRepository using MySQL
type UserRepository struct {
	db     db.DB
	logger log.Logger
}

func NewUserRepository(db db.DB, logger log.Logger) repository.UserRepository {
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
		return domain.NewInvalidUserInputError("Failed to process user preferences", err)
	}

	// Use Execute instead of ExecContext
	_, err = r.db.Execute(
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
				return domain.NewUserAlreadyExistsError(user.Username)
			}
			if strings.Contains(err.Error(), "users.email") {
				return domain.NewUserAlreadyExistsError(user.Email)
			}
			return domain.NewUserAlreadyExistsError(user.ID)
		}
		return domain.Wrap(err, "Failed to create user in database")
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

	// Use QueryRow instead of QueryRowContext
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
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
			return nil, domain.NewUserNotFoundError(id)
		}
		return nil, domain.Wrap(err, "Failed to get user from database")
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	if len(preferencesJSON) > 0 {
		if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
			return nil, domain.NewInvalidUserInputError("Failed to process user preferences", err)
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
    
    // Use QueryRow instead of QueryRowContext
    row := r.db.QueryRow(ctx, query, email)
    err := row.Scan(
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
            return nil, domain.NewUserNotFoundError(email)
        }
        return nil, domain.Wrap(err, "Failed to get user by email from database")
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
            return nil, domain.NewInvalidUserInputError("Failed to process user preferences", err)
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
    
    // Use QueryRow instead of QueryRowContext
    row := r.db.QueryRow(ctx, query, username)
    err := row.Scan(
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
            return nil, domain.NewUserNotFoundError(username)
        }
        return nil, domain.Wrap(err, "Failed to get user by username from database")
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
            return nil, domain.NewInvalidUserInputError("Failed to process user preferences", err)
        }
    }
    
    return user, nil
}

func (r * UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, first_name = ?, last_name = ?, phone = ?, 
		    role = ?, status = ?, email_verified = ?, preferences = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	preferencesJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return domain.NewInvalidUserInputError("Failed to process user preferences", err)
	}

	result, err := r.db.Execute(
		ctx,
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		string(user.Role),
		string(user.Status),
		user.EmailVerified,
		preferencesJSON,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "users.username") {
				return domain.NewUserAlreadyExistsError(user.Username)
			}
			if strings.Contains(err.Error(), "users.email") {
				return domain.NewUserAlreadyExistsError(user.Email)
			}
		}
		return domain.Wrap(err, "Failed to update user in database")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(user.ID)
	}

	return nil
}