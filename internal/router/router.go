package router

import "cinemasys/internal/handler"

type Router struct {
	userHandler        *handler.UserHandler
	ticketHandler      *handler.TicketHandler
	auditoriumHandler  *handler.AuditoriumHandler
	projectionHandler  *handler.ProjectionHandler
	reservationHandler *handler.ReservationHandler
	movieHandler       *handler.MovieHandler
}
