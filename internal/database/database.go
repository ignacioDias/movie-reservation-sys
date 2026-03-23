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

var reservationsCreation = `CREATE TABLE reservations (
    reservation_id BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    projection_id  BIGINT NOT NULL REFERENCES projections(projection_id)
);`

var reservationTicketsCreation = `CREATE TABLE reservation_tickets (
    reservation_ticket_id BIGSERIAL PRIMARY KEY,
    reservation_id        BIGINT NOT NULL REFERENCES reservations(reservation_id) ON DELETE CASCADE,
    ticket_id             BIGINT NOT NULL REFERENCES tickets(ticket_id)
);`

var reservationSeatsCreation = `CREATE TABLE reservation_seats (
    reservation_id BIGINT NOT NULL REFERENCES reservations(reservation_id) ON DELETE CASCADE,
    row            INT NOT NULL,
    col            INT NOT NULL,
    PRIMARY KEY (reservation_id, row, col)
);`

var ticketsCreation = `CREATE TABLE tickets (
    ticket_id   BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    price       NUMERIC(10,2) NOT NULL,
    cant_seats  INT NOT NULL
);`
