package domain

type Reservation struct {
	UserID       int64
	ProjectionID int64
	Amount       float64
	SeatRow      int
	SeatCol      int
}
