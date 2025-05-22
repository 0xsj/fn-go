package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/0xsj/fn-go/pkg/common/errors"
	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
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
	
	metadataJSON, err := json.Marshal(token.Metadata)
	if err != nil {
		return errors.NewInternalError("failed to marshal metadata", err)
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
		return errors.NewDatabaseError("failed to create token", err)
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
			return nil, errors.NewNotFoundError("token not found", err)
		}
		return nil, errors.NewDatabaseError("failed to get token", err)
	}
	
	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}
	
	// Deserialize metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &token.Metadata); err != nil {
			return nil, errors.NewInternalError("failed to unmarshal metadata", err)
		}
	}
	
	return token, nil
}

func (r *TokenRepository) RevokeToken(ctx context.Context, tokenID string) error {
	query := `UPDATE tokens SET revoked_at = ? WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, time.Now(), tokenID)
	if err != nil {
		return errors.NewDatabaseError("failed to revoke token", err)
	}
	
	return nil
}


func (r *TokenRepository) RevokeAllTokensForUser(ctx context.Context, userID string, tokenType string) error {
	var query string
	var args []any
	
	if tokenType != "" {
		query = `UPDATE tokens SET revoked_at = ? WHERE user_id = ? AND type = ? AND revoked_at IS NULL`
		args = []any{time.Now(), userID, tokenType}
	} else {
		query = `UPDATE tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`
		args = []any{time.Now(), userID}
	}
	
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.NewDatabaseError("failed to revoke tokens", err)
	}
	
	return nil
}

func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context) (int, error) {
	query := `DELETE FROM tokens WHERE expires_at < ? OR (revoked_at IS NOT NULL AND revoked_at < ?)`
	
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, now)
	if err != nil {
		return 0, errors.NewDatabaseError("failed to delete expired tokens", err)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.NewDatabaseError("failed to get affected rows", err)
	}
	
	return int(affected), nil
}