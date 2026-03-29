package router

import (
	"cinemasys/internal/cache"
	"cinemasys/internal/database"
	"cinemasys/internal/handler"
	"cinemasys/internal/middleware"
	"net/http"
)

type Router struct {
	mux                *http.ServeMux
	userHandler        *handler.UserHandler
	ticketHandler      *handler.TicketHandler
	auditoriumHandler  *handler.AuditoriumHandler
	projectionHandler  *handler.ProjectionHandler
	reservationHandler *handler.ReservationHandler
	movieHandler       *handler.MovieHandler
	authenticationMw   *middleware.AuthenticationMiddleware
	authorizationMw    *middleware.AuthorizationMiddleware
	rateLimit          *middleware.RateLimitMiddleware
}

func NewRouter(db *database.Database, cache *cache.Cache) *Router {
	// redisClient := database.NewRedisClient(os.Getenv("REDIS_URL"))

	return &Router{
		mux:                http.NewServeMux(),
		userHandler:        handler.NewUserHandler(db.UserRepo, db.SessionRepo),
		ticketHandler:      handler.NewTicketHandler(db.TicketRepo, cache),
		auditoriumHandler:  handler.NewAuditoriumHandler(db.AuditoriumRepo),
		projectionHandler:  handler.NewProjectionHandler(db.ProjectionRepo),
		reservationHandler: handler.NewReservationHandler(db.ReservationRepo),
		movieHandler:       handler.NewMovieHandler(db.MovieRepo, cache),
		authenticationMw:   middleware.NewAuthenticationMiddleware(db.SessionRepo),
		authorizationMw:    middleware.NewAuthorizationMiddleware(db.UserRepo),
		rateLimit:          middleware.NewRateLimitMiddleware(),
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	//front end routes (all in one repo for simplicity, and it's just educational project)
	r.mux.Handle("/", http.FileServer(http.Dir("web")))
	r.mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	r.mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/login.html")
	})
	r.mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/register.html")
	})
	r.mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/profile.html")
	})
	r.mux.HandleFunc("GET /movies/soon", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/future_movies.html")
	})
	r.mux.HandleFunc("GET /movies/id/{movie_id}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/movie.html")
	})
	//session
	r.mux.HandleFunc("POST /api/v1/auth/register", r.rateLimit.RateLimit(r.userHandler.RegisterUser))
	r.mux.HandleFunc("POST /api/v1/auth/login", r.rateLimit.RateLimit(r.userHandler.LoginUser))
	r.mux.HandleFunc("DELETE /api/v1/auth/logout", r.authenticationMw.AuthenticationMiddleware(r.userHandler.LogoutUser))

	//user
	r.mux.HandleFunc("GET /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.GetCurrentUser))
	r.mux.HandleFunc("PUT /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.UpdateUser))
	r.mux.HandleFunc("PUT /api/v1/users/{user_id}/admin", r.authCheck((r.userHandler.MakeUserAdmin)))
	r.mux.HandleFunc("DELETE /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.DeleteMe))
	r.mux.HandleFunc("DELETE /api/v1/users/{user_id}", r.authCheck(r.userHandler.DeleteUser))

	//ticket
	r.mux.HandleFunc("POST /api/v1/tickets", r.authCheck(r.ticketHandler.CreateTicket))
	r.mux.HandleFunc("GET /api/v1/tickets", r.authenticationMw.AuthenticationMiddleware(r.ticketHandler.GetAllTickets))
	r.mux.HandleFunc("PUT /api/v1/tickets/{ticket_id}", r.authCheck(r.ticketHandler.UpdateTicket))
	r.mux.HandleFunc("DELETE /api/v1/tickets/{ticket_id}", r.authCheck(r.ticketHandler.DeleteTicket))

	//auditorium
	r.mux.HandleFunc("GET /api/v1/auditoriums", r.authCheck(r.auditoriumHandler.GetAuditoriums))
	r.mux.HandleFunc("POST /api/v1/auditoriums", r.authCheck(r.auditoriumHandler.CreateAuditorium))
	r.mux.HandleFunc("GET /api/v1/auditoriums/{auditorium_id}", r.authCheck(r.auditoriumHandler.GetAuditoriumByID))
	r.mux.HandleFunc("PUT /api/v1/auditoriums/{auditorium_id}", r.authCheck(r.auditoriumHandler.UpdateAuditorium))
	r.mux.HandleFunc("DELETE /api/v1/auditoriums/{auditorium_id}", r.authCheck(r.auditoriumHandler.DeleteAuditorium))

	//movie
	r.mux.HandleFunc("GET /api/v1/movies", r.movieHandler.GetAllMovies)
	r.mux.HandleFunc("GET /api/v1/movies/soon", r.movieHandler.GetFutureMovies)
	r.mux.HandleFunc("GET /api/v1/movies/available_now", r.movieHandler.GetMoviesWithProjections)

	r.mux.HandleFunc("POST /api/v1/movies", r.authCheck(r.movieHandler.CreateMovie))
	r.mux.HandleFunc("GET /api/v1/movies/id/{movie_id}", r.movieHandler.GetMovieByID)
	r.mux.HandleFunc("PUT /api/v1/movies/id/{movie_id}", r.authCheck(r.movieHandler.UpdateMovie))
	r.mux.HandleFunc("DELETE /api/v1/movies/id/{movie_id}", r.authCheck(r.movieHandler.DeleteMovie))

	//projection
	r.mux.HandleFunc("GET /api/v1/movies/id/{movie_id}/projections", r.projectionHandler.GetAllProjectionsPerMovie)

	r.mux.HandleFunc("POST /api/v1/projections", r.authCheck(r.projectionHandler.CreateProjection))
	r.mux.HandleFunc("GET /api/v1/projections/{projection_id}", r.authCheck(r.projectionHandler.GetProjection))
	r.mux.HandleFunc("PUT /api/v1/projections/{projection_id}", r.authCheck(r.projectionHandler.UpdateProjection))
	r.mux.HandleFunc("DELETE /api/v1/projections/{projection_id}", r.authCheck(r.projectionHandler.DeleteProjection))

	//reservation
	r.mux.HandleFunc("POST /api/v1/reservations", r.authenticationMw.AuthenticationMiddleware(r.reservationHandler.CreateReservation))
	r.mux.HandleFunc("GET /api/v1/users/me/reservations", r.authenticationMw.AuthenticationMiddleware(r.reservationHandler.GetReservationsFromUser))

	return r.mux
}

func (r *Router) authCheck(next http.HandlerFunc) http.HandlerFunc {
	return r.authenticationMw.AuthenticationMiddleware(r.authorizationMw.AuthorizationMiddleware(next))
}
