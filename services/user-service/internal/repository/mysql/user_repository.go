// services/user-service/internal/repository/mysql/user_repository.go
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/user-service/internal/dto"
	"github.com/0xsj/fn-go/services/user-service/internal/repository"

	"github.com/google/uuid"
)

// UserRepository implements repository.UserRepository using MySQL
type UserRepository struct {
	db     *sql.DB
	logger log.Logger
}

// NewUserRepository creates a new MySQL user repository
func NewUserRepository(db *sql.DB, logger log.Logger) repository.UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger.WithLayer("mysql-user-repository"),
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	r.logger.With("username", user.Username).Info("Creating new user")
	
	// Generate a new UUID if ID is empty
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set default values if not provided
	if user.Status == "" {
		user.Status = models.UserStatusPending
	}
	
	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the user
	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (
			id, username, email, password, first_name, last_name, phone, 
			role, status, email_verified, failed_logins, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.Phone,
		user.Role, user.Status, user.EmailVerified, user.FailedLogins, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to insert user")
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Insert user preferences if provided
	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_preferences (
			user_id, theme, notifications_enabled, language, timezone
		) VALUES (?, ?, ?, ?, ?)`,
		user.ID, user.Preferences.Theme, user.Preferences.NotificationsEnabled, 
		user.Preferences.Language, user.Preferences.Timezone,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to insert user preferences")
		return fmt.Errorf("failed to insert user preferences: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		r.logger.With("error", err.Error()).Error("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.With("id", user.ID).Info("User created successfully")
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.logger.With("id", id).Info("Getting user by ID")

	var user models.User
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx,
		`SELECT u.id, u.username, u.email, u.password, u.first_name, u.last_name, u.phone,
			u.role, u.status, u.email_verified, u.failed_logins, u.last_login_at, 
			p.theme, p.notifications_enabled, p.language, p.timezone,
			u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_preferences p ON u.id = p.user_id
		WHERE u.id = ? AND u.deleted_at IS NULL`,
		id,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Phone,
		&user.Role, &user.Status, &user.EmailVerified, &user.FailedLogins, &lastLoginAt,
		&user.Preferences.Theme, &user.Preferences.NotificationsEnabled, &user.Preferences.Language, &user.Preferences.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.With("id", id).Warn("User not found")
			return nil, fmt.Errorf("user not found: %w", err)
		}
		r.logger.With("error", err.Error()).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Set LastLoginAt if not null
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	r.logger.With("id", user.ID).Info("User retrieved successfully")
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.logger.With("email", email).Info("Getting user by email")

	var user models.User
	var lastLoginAt sql.NullTime

	// Query the user
	err := r.db.QueryRowContext(ctx,
		`SELECT u.id, u.username, u.email, u.password, u.first_name, u.last_name, u.phone,
			u.role, u.status, u.email_verified, u.failed_logins, u.last_login_at, 
			p.theme, p.notifications_enabled, p.language, p.timezone,
			u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_preferences p ON u.id = p.user_id
		WHERE u.email = ? AND u.deleted_at IS NULL`,
		email,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Phone,
		&user.Role, &user.Status, &user.EmailVerified, &user.FailedLogins, &lastLoginAt,
		&user.Preferences.Theme, &user.Preferences.NotificationsEnabled, &user.Preferences.Language, &user.Preferences.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.With("email", email).Warn("User not found")
			return nil, fmt.Errorf("user not found: %w", err)
		}
		r.logger.With("error", err.Error()).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Set LastLoginAt if not null
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	r.logger.With("id", user.ID).Info("User retrieved successfully")
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.logger.With("username", username).Info("Getting user by username")

	var user models.User
	var lastLoginAt sql.NullTime

	// Query the user
	err := r.db.QueryRowContext(ctx,
		`SELECT u.id, u.username, u.email, u.password, u.first_name, u.last_name, u.phone,
			u.role, u.status, u.email_verified, u.failed_logins, u.last_login_at, 
			p.theme, p.notifications_enabled, p.language, p.timezone,
			u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_preferences p ON u.id = p.user_id
		WHERE u.username = ? AND u.deleted_at IS NULL`,
		username,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Phone,
		&user.Role, &user.Status, &user.EmailVerified, &user.FailedLogins, &lastLoginAt,
		&user.Preferences.Theme, &user.Preferences.NotificationsEnabled, &user.Preferences.Language, &user.Preferences.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.With("username", username).Warn("User not found")
			return nil, fmt.Errorf("user not found: %w", err)
		}
		r.logger.With("error", err.Error()).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Set LastLoginAt if not null
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	r.logger.With("id", user.ID).Info("User retrieved successfully")
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	r.logger.With("id", user.ID).Info("Updating user")

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update the user
	result, err := tx.ExecContext(ctx,
		`UPDATE users SET
			username = ?, email = ?, first_name = ?, last_name = ?, phone = ?,
			role = ?, status = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		user.Username, user.Email, user.FirstName, user.LastName, user.Phone,
		user.Role, user.Status, user.UpdatedAt, user.ID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", user.ID).Warn("No user found to update")
		return fmt.Errorf("no user found to update")
	}

	// Update user preferences
	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_preferences (
			user_id, theme, notifications_enabled, language, timezone
		) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			theme = VALUES(theme),
			notifications_enabled = VALUES(notifications_enabled),
			language = VALUES(language),
			timezone = VALUES(timezone)`,
		user.ID, user.Preferences.Theme, user.Preferences.NotificationsEnabled,
		user.Preferences.Language, user.Preferences.Timezone,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to update user preferences")
		return fmt.Errorf("failed to update user preferences: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		r.logger.With("error", err.Error()).Error("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.With("id", user.ID).Info("User updated successfully")
	return nil
}

// Delete soft-deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.logger.With("id", id).Info("Deleting user")

	// Soft delete by setting deleted_at timestamp
	now := time.Now()
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL",
		now, id,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", id).Warn("No user found to delete")
		return fmt.Errorf("no user found to delete")
	}

	r.logger.With("id", id).Info("User deleted successfully")
	return nil
}

// List retrieves users based on filter criteria
func (r *UserRepository) List(ctx context.Context, filter dto.ListUsersRequest) ([]*models.User, int, error) {
	r.logger.Info("Listing users with filter")

	// Build the query with filters
	whereClauses := []string{"u.deleted_at IS NULL"}
	args := make([]interface{}, 0)

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		whereClauses = append(whereClauses, "(u.username LIKE ? OR u.email LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ?)")
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Add role filter if present (assuming your DTO has a Role field of type string)
	if roleStr, ok := filter.SortBy, filter.SortBy == "role"; ok && roleStr != "" {
		whereClauses = append(whereClauses, "u.role = ?")
		args = append(args, roleStr)
	}

	// Add status filter if present (assuming your DTO has a Status field of type string)
	if statusStr, ok := filter.SortBy, filter.SortBy == "status"; ok && statusStr != "" {
		whereClauses = append(whereClauses, "u.status = ?")
		args = append(args, statusStr)
	}

	// Construct the WHERE clause
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM users u
		LEFT JOIN user_preferences p ON u.id = p.user_id
		` + whereClause
	
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to count users")
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Add sorting
	orderClause := "ORDER BY u.created_at DESC" // Default sorting
	if filter.SortBy != "" {
		// Map API field names to database column names
		sortFieldMap := map[string]string{
			"id":        "u.id",
			"username":  "u.username",
			"email":     "u.email",
			"firstName": "u.first_name",
			"lastName":  "u.last_name",
			"createdAt": "u.created_at",
		}

		dbField, exists := sortFieldMap[filter.SortBy]
		if exists {
			sortOrder := "ASC"
			if filter.SortOrder == "desc" {
				sortOrder = "DESC"
			}
			orderClause = fmt.Sprintf("ORDER BY %s %s", dbField, sortOrder)
		}
	}

	// Add pagination
	pageSize := 10 // Default page size
	if filter.PageSize > 0 {
		pageSize = filter.PageSize
	}

	page := 1 // Default page
	if filter.Page > 0 {
		page = filter.Page
	}

	offset := (page - 1) * pageSize
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)

	// Final select query
	selectQuery := `
		SELECT u.id, u.username, u.email, u.password, u.first_name, u.last_name, u.phone,
			u.role, u.status, u.email_verified, u.failed_logins, u.last_login_at, 
			p.theme, p.notifications_enabled, p.language, p.timezone,
			u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_preferences p ON u.id = p.user_id
		` + whereClause + ` ` + orderClause + ` ` + limitClause

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to query users")
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	// Process the results
	users := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		var lastLoginAt sql.NullTime
		var deletedAt sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Phone,
			&user.Role, &user.Status, &user.EmailVerified, &user.FailedLogins, &lastLoginAt,
			&user.Preferences.Theme, &user.Preferences.NotificationsEnabled, &user.Preferences.Language, &user.Preferences.Timezone,
			&user.CreatedAt, &user.UpdatedAt, &deletedAt,
		)
		if err != nil {
			r.logger.With("error", err.Error()).Error("Failed to scan user row")
			return nil, 0, fmt.Errorf("failed to scan user row: %w", err)
		}

		// Set nullable fields
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}
		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}

		users = append(users, &user)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		r.logger.With("error", err.Error()).Error("Error iterating user rows")
		return nil, 0, fmt.Errorf("error iterating user rows: %w", err)
	}

	r.logger.With("count", len(users)).With("total", total).Info("Users listed successfully")
	return users, total, nil
}

// AssignRole assigns a role to a user
func (r *UserRepository) AssignRole(ctx context.Context, userID string, role string) error {
	r.logger.With("id", userID).With("role", role).Info("Assigning role to user")

	// Update the user's role
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET role = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		role, time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to assign role")
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to assign role")
		return fmt.Errorf("no user found to assign role")
	}

	r.logger.With("id", userID).With("role", role).Info("Role assigned successfully")
	return nil
}

// RemoveRole removes a role from a user
func (r *UserRepository) RemoveRole(ctx context.Context, userID string, role string) error {
	r.logger.With("id", userID).With("role", role).Info("Removing role from user")

	// First check if the user has the specified role
	var currentRole string
	err := r.db.QueryRowContext(ctx,
		"SELECT role FROM users WHERE id = ? AND deleted_at IS NULL",
		userID,
	).Scan(&currentRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.With("id", userID).Warn("User not found")
			return fmt.Errorf("user not found: %w", err)
		}
		r.logger.With("error", err.Error()).Error("Failed to get user role")
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Check if the user has the role to be removed
	if currentRole != role {
		r.logger.With("id", userID).With("current_role", currentRole).With("role_to_remove", role).Warn("User does not have the specified role")
		return nil // No error, user doesn't have this role
	}

	// Set to a default role (in this case, customer)
	defaultRole := "customer"
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET role = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		defaultRole, time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to remove role")
		return fmt.Errorf("failed to remove role: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to remove role")
		return fmt.Errorf("no user found to remove role")
	}

	r.logger.With("id", userID).With("role", role).With("default_role", defaultRole).Info("Role removed successfully")
	return nil
}

// GetRoles gets the roles for a user
func (r *UserRepository) GetRoles(ctx context.Context, userID string) ([]string, error) {
	r.logger.With("id", userID).Info("Getting roles for user")

	var role string
	err := r.db.QueryRowContext(ctx,
		"SELECT role FROM users WHERE id = ? AND deleted_at IS NULL",
		userID,
	).Scan(&role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.With("id", userID).Warn("User not found")
			return nil, fmt.Errorf("user not found: %w", err)
		}
		r.logger.With("error", err.Error()).Error("Failed to get user role")
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	// Return a slice with a single role for compatibility
	roles := []string{role}
	r.logger.With("id", userID).With("roles", strings.Join(roles, ",")).Info("Roles retrieved successfully")
	return roles, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	r.logger.With("id", userID).Info("Updating user password")

	// Update the password
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET password = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		hashedPassword, time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to update password")
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to update password")
		return fmt.Errorf("no user found to update password")
	}

	r.logger.With("id", userID).Info("Password updated successfully")
	return nil
}

// UpdateLastLoginAt updates a user's last login timestamp
func (r *UserRepository) UpdateLastLoginAt(ctx context.Context, userID string, loginTime time.Time) error {
	r.logger.With("id", userID).Info("Updating user last login timestamp")

	// Update the last login timestamp
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET last_login_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		loginTime, time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to update last login timestamp")
		return fmt.Errorf("failed to update last login timestamp: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to update last login timestamp")
		return fmt.Errorf("no user found to update last login timestamp")
	}

	r.logger.With("id", userID).Info("Last login timestamp updated successfully")
	return nil
}

// IncrementFailedLogins increments a user's failed login count
func (r *UserRepository) IncrementFailedLogins(ctx context.Context, userID string) error {
	r.logger.With("id", userID).Info("Incrementing user failed login count")

	// Increment the failed login count
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET failed_logins = failed_logins + 1, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to increment failed login count")
		return fmt.Errorf("failed to increment failed login count: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to increment failed login count")
		return fmt.Errorf("no user found to increment failed login count")
	}

	r.logger.With("id", userID).Info("Failed login count incremented successfully")
	return nil
}

// ResetFailedLogins resets a user's failed login count to zero
func (r *UserRepository) ResetFailedLogins(ctx context.Context, userID string) error {
	r.logger.With("id", userID).Info("Resetting user failed login count")

	// Reset the failed login count
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET failed_logins = 0, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to reset failed login count")
		return fmt.Errorf("failed to reset failed login count: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to reset failed login count")
		return fmt.Errorf("no user found to reset failed login count")
	}

	r.logger.With("id", userID).Info("Failed login count reset successfully")
	return nil
}

// SetEmailVerified sets a user's email verified status
func (r *UserRepository) SetEmailVerified(ctx context.Context, userID string, verified bool) error {
	r.logger.With("id", userID).With("verified", verified).Info("Setting user email verified status")

	// Set email verified status
	result, err := r.db.ExecContext(ctx,
		"UPDATE users SET email_verified = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		verified, time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to set email verified status")
		return fmt.Errorf("failed to set email verified status: %w", err)
	}

	// Check if the user exists (affected rows)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to get affected rows")
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.With("id", userID).Warn("No user found to set email verified status")
		return fmt.Errorf("no user found to set email verified status")
	}

	r.logger.With("id", userID).With("verified", verified).Info("Email verified status set successfully")
	return nil
}

// UpdatePreferences updates a user's preferences
func (r *UserRepository) UpdatePreferences(ctx context.Context, userID string, preferences models.UserPreferences) error {
	r.logger.With("id", userID).Info("Updating user preferences")

	// Try to update preferences
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_preferences (
			user_id, theme, notifications_enabled, language, timezone
		) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			theme = VALUES(theme),
			notifications_enabled = VALUES(notifications_enabled),
			language = VALUES(language),
			timezone = VALUES(timezone)`,
		userID, preferences.Theme, preferences.NotificationsEnabled,
		preferences.Language, preferences.Timezone,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to update user preferences")
		return fmt.Errorf("failed to update user preferences: %w", err)
	}

	// Check if user exists (query separately)
	var exists bool
	err = r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND deleted_at IS NULL)",
		userID,
	).Scan(&exists)

	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to check if user exists")
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	if !exists {
		r.logger.With("id", userID).Warn("No user found to update preferences")
		return fmt.Errorf("no user found to update preferences")
	}

	// Update the user's updated_at timestamp
	_, err = r.db.ExecContext(ctx,
		"UPDATE users SET updated_at = ? WHERE id = ?",
		time.Now(), userID,
	)
	if err != nil {
		r.logger.With("error", err.Error()).Error("Failed to update user timestamp")
		return fmt.Errorf("failed to update user timestamp: %w", err)
	}

	r.logger.With("id", userID).Info("User preferences updated successfully")
	return nil
}