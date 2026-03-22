package domain

import "errors"

type Reservation struct {
	ReservationID int64    `db:"reservation_id" json:"reservationId"`
	UserID        int64    `db:"user_id" json:"userId"`
	ProjectionID  int64    `db:"projection_id" json:"projectionId"`
	Seats         []Seat   `db:"seats" json:"seats"`
	Tickets       []Ticket `db:"tickets" json:"tickets"`
}

type Seat struct {
	Row int `db:"row" json:"row"`
	Col int `db:"col" json:"col"`
}

type Ticket struct {
	TicketID  int64   `db:"ticket_id" json:"ticketId"`
	Name      string  `db:"name" json:"name"`
	Price     float64 `db:"price" json:"price"`
	CantSeats int     `db:"cant_seats" json:"cant_seats"`
}

func NewReservation(userID, projectionID int64, seats []Seat, tickets []Ticket) (*Reservation, error) {
	if len(seats) != actualCantSeats(tickets) {
		return nil, errors.New("Invalid seats quantity")
	}
	return &Reservation{
		UserID:       userID,
		ProjectionID: projectionID,
		Seats:        seats,
		Tickets:      tickets,
	}, nil
}

func actualCantSeats(tickets []Ticket) int {
	seats := 0
	for _, ticket := range tickets {
		seats += ticket.CantSeats
	}
	return seats
}
