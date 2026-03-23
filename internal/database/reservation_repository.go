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
		query := `INSERT INTO reservation_seats (reservation_id, row, col) VALUES ($1, $2, $3)`
		if _, err := tx.ExecContext(ctx, query, reservation.ReservationID, seat.Row, seat.Col); err != nil {
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

func (rr *ReservationRepository) GetAllUnavailableSeatsFromProjection(ctx context.Context, projectionID int64) ([]domain.Seat, error) {
	query := `SELECT rs.row, rs.col 
          FROM reservation_seats rs 
          INNER JOIN reservations r ON rs.reservation_id = r.reservation_id 
          WHERE r.projection_id = $1`
	var seats []domain.Seat
	if err := rr.db.SelectContext(ctx, &seats, query, projectionID); err != nil {
		return nil, err
	}
	return seats, nil
}

func (rr *ReservationRepository) UpdateReservationSeat(ctx context.Context, reservationID int64, newSeat *domain.Seat, oldSeat *domain.Seat) error {
	query := `UPDATE reservation_seats SET row = $1, col = $2 WHERE reservation_id = $3 AND row = $4 AND col = $5`
	result, err := rr.db.ExecContext(ctx, query, newSeat.Row, newSeat.Col, reservationID, oldSeat.Row, oldSeat.Col)
	return CheckErrResult(result, err, ErrReservationNotFound)
}

func (rr *ReservationRepository) DeleteReservation(ctx context.Context, reservationID int64) error {
	query := `DELETE FROM reservations WHERE reservation_id = $1`
	result, err := rr.db.ExecContext(ctx, query, reservationID)
	return CheckErrResult(result, err, ErrReservationNotFound)
}
