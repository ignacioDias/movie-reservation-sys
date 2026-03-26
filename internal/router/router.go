package router

import (
	"cinemasys/internal/database"
	"cinemasys/internal/handler"
)

type Router struct {
	userHandler        *handler.UserHandler
	ticketHandler      *handler.TicketHandler
	auditoriumHandler  *handler.AuditoriumHandler
	projectionHandler  *handler.ProjectionHandler
	reservationHandler *handler.ReservationHandler
	movieHandler       *handler.MovieHandler
}

func NewRouter(db *database.Database) *Router {
	return &Router{
		userHandler:        handler.NewUserHandler(db.UserRepo, db.SessionRepo),
		ticketHandler:      handler.NewTicketHandler(db.TicketRepo),
		auditoriumHandler:  handler.NewAuditoriumHandler(db.AuditoriumRepo),
		projectionHandler:  handler.NewProjectionHandler(db.ProjectionRepo),
		reservationHandler: handler.NewReservationHandler(db.ReservationRepo),
		movieHandler:       handler.NewMovieHandler(db.MovieRepo),
	}
}
