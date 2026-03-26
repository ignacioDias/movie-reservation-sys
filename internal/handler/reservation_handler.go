package handler

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"cinemasys/internal/middleware"
	"encoding/json"
	"log"
	"net/http"
)

type ReservationHandler struct {
	reservationRepo *database.ReservationRepository
}

type ReservationRequest struct {
	ProjectionID int64           `json:"projectionId"`
	Seats        []domain.Seat   `json:"seats"`
	Tickets      []domain.Ticket `json:"tickets"`
}

func (rh *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservationReq ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&reservationReq); err != nil {
		http.Error(w, "Wrong format for reservation", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	reservation, err := domain.NewReservation(userID, reservationReq.ProjectionID, reservationReq.Seats, reservationReq.Tickets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rh.reservationRepo.CreateReservation(r.Context(), reservation); err != nil {
		http.Error(w, "error while creating reservation", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reservation); err != nil {
		log.Printf("CreateReservation: failed to encode response: %v", err)
	}
}

func (rh *ReservationHandler) GetReservationsFromUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	reservations, err := rh.reservationRepo.GetReservationsPerUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "error getting reservations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reservations); err != nil {
		log.Printf("GetReservations: failed to encode response: %v", err)
	}
}
