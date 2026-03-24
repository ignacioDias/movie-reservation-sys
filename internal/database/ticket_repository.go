package database

import (
	"cinemasys/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

var ErrTicketNotFound = errors.New("Ticket not found")

func (tr *TicketRepository) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	query := `INSERT INTO tickets (name, price, cant_seats) VALUES ($1, $2, $3) RETURNING ticket_id`
	return tr.db.QueryRowContext(ctx, query, ticket.Name, ticket.Price, ticket.CantSeats).Scan(&ticket.TicketID)
}

func (tr *TicketRepository) GetTicketByID(ctx context.Context, ticketID int64) (*domain.Ticket, error) {
	query := `SELECT * FROM tickets WHERE ticket_id = $1`
	var ticket domain.Ticket
	if err := tr.db.GetContext(ctx, &ticket, query, ticketID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

func (tr *TicketRepository) GetAllTickets(ctx context.Context) ([]domain.Ticket, error) {
	query := `SELECT * FROM tickets`
	var tickets []domain.Ticket
	if err := tr.db.SelectContext(ctx, &tickets, query); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (tr *TicketRepository) UpdateTicket(ctx context.Context, ticket *domain.Ticket) error {
	query := `UPDATE tickets SET name = $1, price = $2, cant_seats = $3 WHERE ticket_id = $4`
	result, err := tr.db.ExecContext(ctx, query, ticket.Name, ticket.Price, ticket.CantSeats, ticket.TicketID)
	return CheckErrResult(result, err, ErrTicketNotFound)
}

func (tr *TicketRepository) DeleteTicket(ctx context.Context, ticketID int64) error {
	query := `DELETE FROM tickets WHERE ticket_id = $1`
	result, err := tr.db.ExecContext(ctx, query, ticketID)
	return CheckErrResult(result, err, ErrTicketNotFound)
}
