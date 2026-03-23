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
	TicketRepo      *TicketRepository
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
		TicketRepo:      NewTicketRepository(db),
	}
}

var usersCreation = `CREATE TABLE IF NOT EXISTS users (
    user_id   BIGSERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	hashed_password TEXT NOT NULL,
	document_number TEXT UNIQUE NOT NULL,
	profile_picture TEXT NOT NULL,
	role INT NOT NULL CHECK (role IN (0, 1)) DEFAULT 0
);`

var reservationsCreation = `CREATE TABLE IF NOT EXISTS reservations (
    reservation_id BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	projection_id  BIGINT NOT NULL REFERENCES projections(projection_id),
	UNIQUE (reservation_id, projection_id)
);`

var reservationTicketsCreation = `CREATE TABLE IF NOT EXISTS reservation_tickets (
    reservation_ticket_id BIGSERIAL PRIMARY KEY,
    reservation_id        BIGINT NOT NULL REFERENCES reservations(reservation_id) ON DELETE CASCADE,
    ticket_id             BIGINT NOT NULL REFERENCES tickets(ticket_id)
);`

var reservationSeatsCreation = `CREATE TABLE IF NOT EXISTS reservation_seats (
	projection_id BIGINT NOT NULL,
	reservation_id BIGINT NOT NULL,
    row            INT NOT NULL,
    col            INT NOT NULL,
	PRIMARY KEY (projection_id, row, col),
	FOREIGN KEY (reservation_id, projection_id) REFERENCES reservations(reservation_id, projection_id) ON DELETE CASCADE
);`

var ticketsCreation = `CREATE TABLE IF NOT EXISTS tickets (
    ticket_id   BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    price       NUMERIC(10,2) NOT NULL,
    cant_seats  INT NOT NULL
);`

var createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);`
