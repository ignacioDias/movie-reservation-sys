package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    int64     `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

func NewSession(userID int64) *Session {
	return &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(30 * time.Minute),
	}
}
