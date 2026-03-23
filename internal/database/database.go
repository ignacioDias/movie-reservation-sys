package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

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

var usersTableCreation = `CREATE TABLE IF NOT EXISTS users (
    user_id   BIGSERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	hashed_password TEXT NOT NULL,
	document_number TEXT UNIQUE NOT NULL,
	profile_picture TEXT NOT NULL,
	role INT NOT NULL CHECK (role IN (0, 1)) DEFAULT 0
);`

var moviesTableCreation = `
CREATE TABLE IF NOT EXISTS movies (
	movie_id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	poster_image_url TEXT NOT NULL,
	trailer_url TEXT NOT NULL,
	genres TEXT[] NOT NULL,
	release_date TIMESTAMPTZ NOT NULL
)
`

var auditoriumsTableCreation = `
CREATE TABLE IF NOT EXISTS auditoriums (
	auditorium_id BIGSERIAL PRIMARY KEY,
    cant_rows     INT NOT NULL CHECK (cant_rows > 0 AND cant_rows <= 100),
    cant_cols     INT NOT NULL CHECK (cant_cols > 0 AND cant_cols <= 100),
	name TEXT NOT NULL
);`

var projectionsTableCreation = `
CREATE TABLE IF NOT EXISTS projections (
	projection_id BIGSERIAL PRIMARY KEY,
	auditorium_id BIGINT NOT NULL REFERENCES auditoriums(auditorium_id) ON DELETE CASCADE,
	movie_id BIGINT NOT NULL REFERENCES movies(movie_id) ON DELETE CASCADE,
	screening_format TEXT NOT NULL CHECK (screening_format IN ('2D', '3D')),
	language TEXT NOT NULL CHECK (language IN ('Spanish', 'Original', 'Other')),
	starts_at TIMESTAMPTZ NOT NULL
);`

var reservationsTableCreation = `CREATE TABLE IF NOT EXISTS reservations (
    reservation_id BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	projection_id  BIGINT NOT NULL REFERENCES projections(projection_id), 
	UNIQUE (reservation_id, projection_id)
);`

var ticketsTableCreation = `CREATE TABLE IF NOT EXISTS tickets (
    ticket_id   BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    price       NUMERIC(10,2) NOT NULL CHECK (price > 0),
    cant_seats  INT NOT NULL CHECK (cant_seats > 0)
);`

var reservationTicketsTableCreation = `CREATE TABLE IF NOT EXISTS reservation_tickets (
    reservation_ticket_id BIGSERIAL PRIMARY KEY,
    reservation_id        BIGINT NOT NULL REFERENCES reservations(reservation_id) ON DELETE CASCADE,
    ticket_id             BIGINT NOT NULL REFERENCES tickets(ticket_id)
);`

var reservationSeatsTableCreation = `CREATE TABLE IF NOT EXISTS reservation_seats (
	projection_id BIGINT NOT NULL,
	reservation_id BIGINT NOT NULL,
    row            INT NOT NULL CHECK (row >= 0 AND row < 100),
    col            INT NOT NULL CHECK (col >= 0 AND col < 100),
	PRIMARY KEY (projection_id, row, col),
	FOREIGN KEY (reservation_id, projection_id) REFERENCES reservations(reservation_id, projection_id) ON DELETE CASCADE
);`

var createTableSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);`

func (d *Database) InitDB() error {
	tx, err := d.DB.Beginx()
	if err != nil {
		return fmt.Errorf("begin init db transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	tableStatements := []struct {
		name string
		ddl  string
	}{
		{name: "users", ddl: usersTableCreation},
		{name: "movies", ddl: moviesTableCreation},
		{name: "auditoriums", ddl: auditoriumsTableCreation},
		{name: "tickets", ddl: ticketsTableCreation},
		{name: "projections", ddl: projectionsTableCreation},
		{name: "reservations", ddl: reservationsTableCreation},
		{name: "reservation_tickets", ddl: reservationTicketsTableCreation},
		{name: "reservation_seats", ddl: reservationSeatsTableCreation},
		{name: "sessions", ddl: createTableSessionsTable},
	}

	for _, tableStmt := range tableStatements {
		if _, err := tx.Exec(tableStmt.ddl); err != nil {
			return fmt.Errorf("create table %s: %w", tableStmt.name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit init db transaction: %w", err)
	}

	return nil
}
