// services/user-service/internal/repository/mysql/user_repository.go
package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0xsj/fn-go/pkg/common/db"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/domain"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
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


func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Soft delete by setting deleted_at timestamp
	query := `
		UPDATE users 
		SET deleted_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.Execute(ctx, query, now, now, id)
	if err != nil {
		return domain.Wrap(err, "Failed to delete user from database")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(id)
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context, req dto.ListUsersRequest) ([]*models.User, int, error) {
	// Build the WHERE clause
	whereClause, args := r.buildWhereClause(req)
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	row := r.db.QueryRow(ctx, countQuery, args...)
	
	var totalCount int
	if err := row.Scan(&totalCount); err != nil {
		return nil, 0, domain.Wrap(err, "Failed to count users")
	}

	// Build the main query
	orderClause := r.buildOrderClause(req)
	limitClause := r.buildLimitClause(req)
	
	query := fmt.Sprintf(`
		SELECT id, username, email, password, first_name, last_name, phone, 
		       role, status, last_login_at, failed_logins, email_verified, 
		       preferences, created_at, updated_at, deleted_at
		FROM users 
		%s %s %s
	`, whereClause, orderClause, limitClause)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, domain.Wrap(err, "Failed to list users from database")
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var preferencesJSON []byte
		var lastLoginAt sql.NullTime
		var deletedAt sql.NullTime

		err := rows.Scan(
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
			return nil, 0, domain.Wrap(err, "Failed to scan user row")
		}

		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}
		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}

		if len(preferencesJSON) > 0 {
			if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
				return nil, 0, domain.NewInvalidUserInputError("Failed to process user preferences", err)
			}
		}

		// Remove password before adding to results
		user.Password = ""
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, domain.Wrap(err, "Error iterating user rows")
	}

	return users, totalCount, nil
}


func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `
		UPDATE users 
		SET password = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.Execute(ctx, query, hashedPassword, time.Now(), userID)
	if err != nil {
		return domain.Wrap(err, "Failed to update user password")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}


func (r *UserRepository) UpdateLastLoginAt(ctx context.Context, userID string, loginTime time.Time) error {
	query := `
		UPDATE users 
		SET last_login_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.Execute(ctx, query, loginTime, time.Now(), userID)
	if err != nil {
		return domain.Wrap(err, "Failed to update user last login time")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

func (r *UserRepository) IncrementFailedLogins(ctx context.Context, userID string) error {
	query := `
		UPDATE users 
		SET failed_logins = failed_logins + 1, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.Execute(ctx, query, time.Now(), userID)
	if err != nil {
		return domain.Wrap(err, "Failed to increment failed logins")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

func (r *UserRepository) ResetFailedLogins(ctx context.Context, userID string) error {
	query := `
		UPDATE users 
		SET failed_logins = 0, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.Execute(ctx, query, time.Now(), userID)
	if err != nil {
		return domain.Wrap(err, "Failed to reset failed logins")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

func (r *UserRepository) SetEmailVerified(ctx context.Context, userID string, verified bool) error {
	query := `
		UPDATE users 
		SET email_verified = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.Execute(ctx, query, verified, time.Now(), userID)
	if err != nil {
		return domain.Wrap(err, "Failed to update email verification status")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}

func (r *UserRepository) UpdatePreferences(ctx context.Context, userID string, preferences models.UserPreferences) error {
	preferencesJSON, err := json.Marshal(preferences)
	if err != nil {
		return domain.NewInvalidUserInputError("Failed to process user preferences", err)
	}

	query := `
		UPDATE users 
		SET preferences = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.Execute(ctx, query, preferencesJSON, time.Now(), userID)
	if err != nil {
		return domain.Wrap(err, "Failed to update user preferences")
	}

	if result == 0 {
		return domain.NewUserNotFoundError(userID)
	}

	return nil
}


func (r *UserRepository) buildWhereClause(req dto.ListUsersRequest) (string, []interface{}) {
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}

	// Search filter
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		conditions = append(conditions, "(username LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ?)")
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Role filter
	if len(req.Roles) > 0 {
		placeholders := make([]string, len(req.Roles))
		for i, role := range req.Roles {
			placeholders[i] = "?"
			args = append(args, role)
		}
		conditions = append(conditions, fmt.Sprintf("role IN (%s)", strings.Join(placeholders, ",")))
	}

	// Active status filter
	if req.IsActive != nil {
		if *req.IsActive {
			conditions = append(conditions, "status = ?")
			args = append(args, "active")
		} else {
			conditions = append(conditions, "status != ?")
			args = append(args, "active")
		}
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (r *UserRepository) buildOrderClause(req dto.ListUsersRequest) string {
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	// Validate sort field to prevent SQL injection
	validSortFields := map[string]bool{
		"id":         true,
		"username":   true,
		"email":      true,
		"first_name": true,
		"last_name":  true,
		"created_at": true,
		"updated_at": true,
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := strings.ToUpper(req.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", sortBy, sortOrder)
}

func (r *UserRepository) buildLimitClause(req dto.ListUsersRequest) string {
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)
}