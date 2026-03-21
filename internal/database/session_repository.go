package database

import (
	"cinemasys/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

var ErrSessionNotFound = errors.New("Session not found")

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (sessRepo *SessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
	INSERT INTO sessions (id, user_id, expires_at)
	VALUES (:id, :user_id, :expires_at)
	RETURNING created_at
	`
	stmt, err := sessRepo.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, session, session)
}

func (sessRepo *SessionRepository) FindSessionByID(ctx context.Context, id string) (*domain.Session, error) {
	var session domain.Session
	query := "SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1 AND expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')"
	if err := sessRepo.db.GetContext(ctx, &session, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (sessRepo *SessionRepository) DeleteSessionByID(ctx context.Context, id string) error {
	query := "DELETE FROM sessions WHERE id = $1"
	result, err := sessRepo.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrSessionNotFound)
}

func (sessRepo *SessionRepository) DeleteSessionsByUserID(ctx context.Context, userID int64) error {
	query := "DELETE FROM sessions WHERE user_id = $1"
	result, err := sessRepo.db.ExecContext(ctx, query, userID)
	return CheckErrResult(result, err, ErrSessionNotFound)
}
