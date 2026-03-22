package database

import "github.com/jmoiron/sqlx"

type Database struct {
	DB              *sqlx.DB
	SessionRepo     *SessionRepository
	UserRepo        *UserRepository
	AuditoriumRepo  *AuditoriumRepository
	MovieRepo       *MovieRepository
	ProjectionRepo  *ProjectionRepository
	ReservationRepo *ReservationRepository
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		DB:              db,
		SessionRepo:     NewSessionRepository(db),
		UserRepo:        NewUserRepository(db),
		AuditoriumRepo:  NewAuditoriumRepository(db),
		MovieRepo:       NewMovieRepository(db),
		ProjectionRepo:  NewProjectionRepository(db),
		ReservationRepo: NewReservationRepository(db),
	}
}
