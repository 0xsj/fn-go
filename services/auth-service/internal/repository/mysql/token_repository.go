// services/auth-service/internal/repository/mysql/token_repository.go
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/repository"
)

type TokenRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewTokenRepository(db *sql.DB, logger log.Logger) repository.TokenRepository {
	return &TokenRepository{
		db:     db,
		logger: logger.WithLayer("mysql-token-repository"),
	}
}

func (r *TokenRepository) CreateToken(ctx context.Context, token *models.Token) error {
	query := `
		INSERT INTO tokens (
			id, user_id, type, value, expires_at, created_at, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	var metadataJSON []byte
	var err error
	if token.Metadata != nil {
		metadataJSON, err = json.Marshal(token.Metadata)
		if err != nil {
			return domain.NewInvalidAuthInputError("failed to marshal metadata", err)
		}
	}
	
	_, err = r.db.ExecContext(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.Type,
		token.Value,
		token.ExpiresAt,
		token.CreatedAt,
		metadataJSON,
	)
	
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to create token in database"),
			"create_token",
		)
	}
	
	return nil
}

func (r *TokenRepository) GetTokenByValue(ctx context.Context, value string) (*models.Token, error) {
	query := `
		SELECT id, user_id, type, value, expires_at, revoked_at, created_at, metadata
		FROM tokens
		WHERE value = ? AND (revoked_at IS NULL OR revoked_at > NOW())
	`
	
	token := &models.Token{}
	var metadataJSON []byte
	var revokedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&token.ID,
		&token.UserID,
		&token.Type,
		&token.Value,
		&token.ExpiresAt,
		&revokedAt,
		&token.CreatedAt,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewTokenNotFoundError("token_value")
		}
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get token from database"),
			"get_token_by_value",
		)
	}
	
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	
	// Deserialize metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &token.Metadata); err != nil {
			return nil, domain.NewInvalidAuthInputError("failed to unmarshal metadata", err)
		}
	}
	
	return token, nil
}

func (r *TokenRepository) GetTokenByID(ctx context.Context, id string) (*models.Token, error) {
	query := `
		SELECT id, user_id, type, value, expires_at, revoked_at, created_at, metadata
		FROM tokens
		WHERE id = ?
	`
	
	token := &models.Token{}
	var metadataJSON []byte
	var revokedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&token.ID,
		&token.UserID,
		&token.Type,
		&token.Value,
		&token.ExpiresAt,
		&revokedAt,
		&token.CreatedAt,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewTokenNotFoundError(id)
		}
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get token by ID from database"),
			"get_token_by_id",
		)
	}
	
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	
	// Deserialize metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &token.Metadata); err != nil {
			return nil, domain.NewInvalidAuthInputError("failed to unmarshal metadata", err)
		}
	}
	
	return token, nil
}

func (r *TokenRepository) GetTokensByUserID(ctx context.Context, userID string, tokenType string) ([]*models.Token, error) {
	var query string
	var args []interface{}
	
	if tokenType != "" {
		query = `
			SELECT id, user_id, type, value, expires_at, revoked_at, created_at, metadata
			FROM tokens
			WHERE user_id = ? AND type = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{userID, tokenType}
	} else {
		query = `
			SELECT id, user_id, type, value, expires_at, revoked_at, created_at, metadata
			FROM tokens
			WHERE user_id = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{userID}
	}
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get tokens by user ID from database"),
			"get_tokens_by_user_id",
		)
	}
	defer rows.Close()
	
	var tokens []*models.Token
	for rows.Next() {
		token := &models.Token{}
		var metadataJSON []byte
		var revokedAt sql.NullTime
		
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Type,
			&token.Value,
			&token.ExpiresAt,
			&revokedAt,
			&token.CreatedAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, domain.WithOperation(
				domain.Wrap(err, "failed to scan token row"),
				"get_tokens_by_user_id",
			)
		}
		
		if revokedAt.Valid {
			token.RevokedAt = &revokedAt.Time
		}
		
		// Deserialize metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &token.Metadata); err != nil {
				r.logger.With("error", err.Error()).Warn("Failed to unmarshal token metadata")
				// Continue without metadata rather than failing
			}
		}
		
		tokens = append(tokens, token)
	}
	
	if err := rows.Err(); err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "error iterating token rows"),
			"get_tokens_by_user_id",
		)
	}
	
	return tokens, nil
}

func (r *TokenRepository) RevokeToken(ctx context.Context, tokenID string) error {
	query := `UPDATE tokens SET revoked_at = ? WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), tokenID)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to revoke token in database"),
			"revoke_token",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"revoke_token",
		)
	}
	
	if affected == 0 {
		return domain.NewTokenNotFoundError(tokenID)
	}
	
	return nil
}

func (r *TokenRepository) RevokeAllTokensForUser(ctx context.Context, userID string, tokenType string) error {
	var query string
	var args []interface{}
	
	if tokenType != "" {
		query = `UPDATE tokens SET revoked_at = ? WHERE user_id = ? AND type = ? AND revoked_at IS NULL`
		args = []interface{}{time.Now(), userID, tokenType}
	} else {
		query = `UPDATE tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`
		args = []interface{}{time.Now(), userID}
	}
	
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to revoke tokens for user"),
			"revoke_all_tokens_for_user",
		)
	}
	
	return nil
}

func (r *TokenRepository) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	query := `SELECT revoked_at FROM tokens WHERE id = ?`
	
	var revokedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, tokenID).Scan(&revokedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, domain.NewTokenNotFoundError(tokenID)
		}
		return false, domain.WithOperation(
			domain.Wrap(err, "failed to check token revocation status"),
			"is_token_revoked",
		)
	}
	
	return revokedAt.Valid, nil
}

func (r *TokenRepository) UpdateToken(ctx context.Context, token *models.Token) error {
	query := `
		UPDATE tokens 
		SET type = ?, value = ?, expires_at = ?, revoked_at = ?, metadata = ?
		WHERE id = ?
	`
	
	var metadataJSON []byte
	var err error
	if token.Metadata != nil {
		metadataJSON, err = json.Marshal(token.Metadata)
		if err != nil {
			return domain.NewInvalidAuthInputError("failed to marshal metadata", err)
		}
	}
	
	result, err := r.db.ExecContext(
		ctx,
		query,
		token.Type,
		token.Value,
		token.ExpiresAt,
		token.RevokedAt,
		metadataJSON,
		token.ID,
	)
	
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to update token in database"),
			"update_token",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"update_token",
		)
	}
	
	if affected == 0 {
		return domain.NewTokenNotFoundError(token.ID)
	}
	
	return nil
}

func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context) (int, error) {
	query := `DELETE FROM tokens WHERE expires_at < ? OR (revoked_at IS NOT NULL AND revoked_at < ?)`
	
	cutoffTime := time.Now().Add(-24 * time.Hour) // Keep revoked tokens for 24 hours for audit
	result, err := r.db.ExecContext(ctx, query, time.Now(), cutoffTime)
	if err != nil {
		return 0, domain.WithOperation(
			domain.Wrap(err, "failed to delete expired tokens"),
			"delete_expired_tokens",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"delete_expired_tokens",
		)
	}
	
	return int(affected), nil
}