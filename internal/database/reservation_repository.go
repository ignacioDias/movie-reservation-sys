package database

import (
	"cinemasys/internal/domain"
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrReservationNotFound = errors.New("reservation not found")

type ReservationRepository struct {
	db *sqlx.DB
}

func NewReservationRepository(db *sqlx.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (rr *ReservationRepository) CreateReservation(ctx context.Context, reservation *domain.Reservation) error {
	tx, err := rr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO reservations (user_id, projection_id) VALUES ($1, $2) RETURNING reservation_id`
	if err := tx.QueryRowContext(ctx, query, reservation.UserID, reservation.ProjectionID).Scan(&reservation.ReservationID); err != nil {
		return err
	}
	for _, seat := range reservation.Seats {
		query := `INSERT INTO reservation_seats (reservation_id, projection_id, row, col) VALUES ($1, $2, $3, $4)`
		if _, err := tx.ExecContext(ctx, query, reservation.ReservationID, reservation.ProjectionID, seat.Row, seat.Col); err != nil {
			return err
		}
	}
	for _, ticket := range reservation.Tickets {
		query := `INSERT INTO reservation_tickets (reservation_id, ticket_id) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, query, reservation.ReservationID, ticket.TicketID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (rr *ReservationRepository) GetReservationsPerUser(ctx context.Context, userID int64) ([]domain.Reservation, error) {
	query := `SELECT * FROM reservations WHERE user_id = $1`
	querySeats := `SELECT row, col FROM reservation_seats WHERE reservation_id = $1`
	queryTickets := `
		SELECT t.ticket_id, t.name, t.price, t.cant_seats
		FROM tickets t
		INNER JOIN reservation_tickets rt ON t.ticket_id = rt.ticket_id
		WHERE rt.reservation_id = $1`

	var reservations []domain.Reservation
	if err := rr.db.SelectContext(ctx, &reservations, query, userID); err != nil {
		return nil, err
	}

	for i := range reservations {
		if err := rr.db.SelectContext(ctx, &reservations[i].Seats, querySeats, reservations[i].ReservationID); err != nil {
			return nil, err
		}
		if err := rr.db.SelectContext(ctx, &reservations[i].Tickets, queryTickets, reservations[i].ReservationID); err != nil {
			return nil, err
		}
	}

	return reservations, nil
}

func (rr *ReservationRepository) GetAllUnavailableSeatsFromProjection(ctx context.Context, projectionID int64) ([]domain.Seat, error) {
	query := `SELECT row, col FROM reservation_seats WHERE projection_id = $1`
	var seats []domain.Seat
	if err := rr.db.SelectContext(ctx, &seats, query, projectionID); err != nil {
		return nil, err
	}
	return seats, nil
}
