package handler

import (
	"cinemasys/internal/cache"
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TicketHandler struct {
	ticketRepo *database.TicketRepository
	cache      *cache.Cache
}

type TicketRequest struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CantSeats int     `json:"cant_seats"`
}
type TicketUpdateRequest struct {
	Name      *string  `json:"name"`
	Price     *float64 `json:"price"`
	CantSeats *int     `json:"cant_seats"`
}

func NewTicketHandler(ticketRepo *database.TicketRepository, cache *cache.Cache) *TicketHandler {
	return &TicketHandler{
		ticketRepo: ticketRepo,
		cache:      cache,
	}
}

func (th *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var tickerReq TicketRequest
	if err := json.NewDecoder(r.Body).Decode(&tickerReq); err != nil {
		http.Error(w, "wrong format ticket request", http.StatusBadRequest)
		return
	}
	ticket, err := domain.NewTicket(tickerReq.Name, tickerReq.Price, tickerReq.CantSeats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := th.ticketRepo.CreateTicket(r.Context(), ticket); err != nil {
		http.Error(w, "error creating ticket", http.StatusInternalServerError)
		return
	}
	if err := th.cache.Delete("tickets"); err != nil {
		log.Printf("CreateTickets: failed to invalidate tickets cache: %v", err)
	}
	WriteResponseWithEncoder(w, ticket, http.StatusCreated)
}

func (th *TicketHandler) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	var tickets []domain.Ticket
	if err := th.cache.Get("tickets", &tickets); err != nil {
		tickets, err = th.ticketRepo.GetAllTickets(r.Context())
		if err != nil {
			http.Error(w, "error getting tickets", http.StatusInternalServerError)
			return
		}
		if err := th.cache.Set("tickets", tickets, time.Hour); err != nil {
			log.Printf("GetAllTickets: failed to cache tickets: %v", err)
		}
	}
	WriteResponseWithEncoder(w, tickets, http.StatusOK)
}

func (th *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	var updateReq TicketUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "error in update request", http.StatusBadRequest)
		return
	}
	id := r.PathValue("ticket_id")
	ticketID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong format ticket id passed", http.StatusBadRequest)
		return
	}
	ticket, err := th.ticketRepo.GetTicketByID(r.Context(), ticketID)
	if err != nil {
		if errors.Is(err, database.ErrTicketNotFound) {
			http.Error(w, "ticket not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error getting ticket", http.StatusInternalServerError)
		return
	}
	if updateReq.CantSeats != nil {
		if *updateReq.CantSeats <= 0 {
			http.Error(w, "seats can't be <= 0", http.StatusBadRequest)
			return
		}
		ticket.CantSeats = *updateReq.CantSeats
	}
	if updateReq.Price != nil {
		if *updateReq.Price < 0 {
			http.Error(w, "price can't be < 0", http.StatusBadRequest)
			return
		}
		ticket.Price = *updateReq.Price
	}
	if updateReq.Name != nil {
		if *updateReq.Name == "" {
			http.Error(w, "name empty", http.StatusBadRequest)
			return
		}
		ticket.Name = *updateReq.Name
	}
	if err := th.ticketRepo.UpdateTicket(r.Context(), ticket); err != nil {
		if errors.Is(err, database.ErrTicketNotFound) {
			http.Error(w, "ticket not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating ticket", http.StatusInternalServerError)
		return
	}
	if err := th.cache.Delete("tickets"); err != nil {
		log.Printf("UpdateTickets: failed to invalidate tickets cache: %v", err)
	}
	WriteResponseWithEncoder(w, ticket, http.StatusOK)
}

func (th *TicketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("ticket_id")
	ticketID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for ticket id passed", http.StatusBadRequest)
		return
	}
	if err := th.ticketRepo.DeleteTicket(r.Context(), ticketID); err != nil {
		if errors.Is(err, database.ErrTicketNotFound) {
			http.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error while deleting ticket", http.StatusInternalServerError)
		return
	}
	if err := th.cache.Delete("tickets"); err != nil {
		log.Printf("DeleteTickets: failed to invalidate tickets cache: %v", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
