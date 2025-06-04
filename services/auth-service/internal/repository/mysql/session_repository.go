// services/auth-service/internal/repository/mysql/session_repository.go
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
	"github.com/0xsj/fn-go/pkg/models"
	"github.com/0xsj/fn-go/services/auth-service/internal/domain"
	"github.com/0xsj/fn-go/services/auth-service/internal/repository"
)

type SessionRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewSessionRepository(db *sql.DB, logger log.Logger) repository.SessionRepository {
	return &SessionRepository{
		db:     db,
		logger: logger.WithLayer("mysql-session-repository"),
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, refresh_token, user_agent, ip_address, 
			last_active, expires_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.LastActive,
		session.ExpiresAt,
		session.CreatedAt,
	)
	
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to create session in database"),
			"create_session",
		)
	}
	
	return nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address,
		       last_active, expires_at, created_at
		FROM sessions
		WHERE id = ? AND expires_at > NOW()
	`
	
	session := &models.Session{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.LastActive,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewSessionNotFoundError(id)
		}
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get session from database"),
			"get_session_by_id",
		)
	}
	
	return session, nil
}

func (r *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address,
		       last_active, expires_at, created_at
		FROM sessions
		WHERE refresh_token = ? AND expires_at > NOW()
	`
	
	session := &models.Session{}
	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.LastActive,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewSessionNotFoundError("refresh_token")
		}
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get session by refresh token"),
			"get_session_by_refresh_token",
		)
	}
	
	return session, nil
}

func (r *SessionRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address,
		       last_active, expires_at, created_at
		FROM sessions
		WHERE user_id = ? AND expires_at > NOW()
		ORDER BY last_active DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "failed to get sessions for user"),
			"get_sessions_by_user_id",
		)
	}
	defer rows.Close()
	
	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RefreshToken,
			&session.UserAgent,
			&session.IPAddress,
			&session.LastActive,
			&session.ExpiresAt,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, domain.WithOperation(
				domain.Wrap(err, "failed to scan session row"),
				"get_sessions_by_user_id",
			)
		}
		sessions = append(sessions, session)
	}
	
	if err := rows.Err(); err != nil {
		return nil, domain.WithOperation(
			domain.Wrap(err, "error iterating session rows"),
			"get_sessions_by_user_id",
		)
	}
	
	return sessions, nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE sessions 
		SET user_agent = ?, ip_address = ?, last_active = ?, expires_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(
		ctx,
		query,
		session.UserAgent,
		session.IPAddress,
		session.LastActive,
		session.ExpiresAt,
		session.ID,
	)
	
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to update session in database"),
			"update_session",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"update_session",
		)
	}
	
	if affected == 0 {
		return domain.NewSessionNotFoundError(session.ID)
	}
	
	return nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to delete session from database"),
			"delete_session",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"delete_session",
		)
	}
	
	if affected == 0 {
		return domain.NewSessionNotFoundError(id)
	}
	
	return nil
}

func (r *SessionRepository) DeleteAllSessionsForUser(ctx context.Context, userID string) (int, error) {
	query := `DELETE FROM sessions WHERE user_id = ?`
	
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, domain.WithOperation(
			domain.Wrap(err, "failed to delete user sessions from database"),
			"delete_all_sessions_for_user",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"delete_all_sessions_for_user",
		)
	}
	
	return int(affected), nil
}

func (r *SessionRepository) UpdateSessionLastActive(ctx context.Context, id string, lastActive time.Time) error {
	query := `UPDATE sessions SET last_active = ? WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, lastActive, id)
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to update session last active"),
			"update_session_last_active",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"update_session_last_active",
		)
	}
	
	if affected == 0 {
		return domain.NewSessionNotFoundError(id)
	}
	
	return nil
}

func (r *SessionRepository) DeleteExpiredSessions(ctx context.Context) (int, error) {
	query := `DELETE FROM sessions WHERE expires_at <= NOW()`
	
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, domain.WithOperation(
			domain.Wrap(err, "failed to delete expired sessions"),
			"delete_expired_sessions",
		)
	}
	
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, domain.WithOperation(
			domain.Wrap(err, "failed to get affected rows"),
			"delete_expired_sessions",
		)
	}
	
	return int(affected), nil
}