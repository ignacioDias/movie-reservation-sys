package handler

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type AuditoriumRequest struct {
	CantRows int    `json:"cantRows"`
	CantCols int    `json:"cantCols"`
	Name     string `json:"name"`
}
type AuditoriumUpdateRequest struct {
	CantRows *int    `json:"cantRows"`
	CantCols *int    `json:"cantCols"`
	Name     *string `json:"name"`
}

type AuditoriumHandler struct {
	auditoriumRepo *database.AuditoriumRepository
}

func NewAuditoriumHandler(repo *database.AuditoriumRepository) *AuditoriumHandler {
	return &AuditoriumHandler{auditoriumRepo: repo}
}

func (ah *AuditoriumHandler) CreateAuditorium(w http.ResponseWriter, r *http.Request) {
	var auditoriumReq AuditoriumRequest
	if err := json.NewDecoder(r.Body).Decode(&auditoriumReq); err != nil {
		http.Error(w, "Wrong format for auditorium", http.StatusBadRequest)
		return
	}
	auditorium := domain.NewAuditorium(auditoriumReq.CantRows, auditoriumReq.CantCols, auditoriumReq.Name)
	if err := ah.auditoriumRepo.CreateAuditorium(r.Context(), auditorium); err != nil {
		http.Error(w, "error while creating auditorium", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, auditorium, http.StatusCreated)
}

func (ah *AuditoriumHandler) GetAuditoriumByID(w http.ResponseWriter, r *http.Request) {
	auditorium, err := ah.getAuditoriumFromPath(r)
	if err != nil {
		writeAuditoriumError(w, err)
		return
	}
	WriteResponseWithEncoder(w, auditorium, http.StatusOK)
}

func (ah *AuditoriumHandler) UpdateAuditorium(w http.ResponseWriter, r *http.Request) {
	var auditoriumUpdateReq AuditoriumUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&auditoriumUpdateReq); err != nil {
		http.Error(w, "wrong format for update", http.StatusBadRequest)
		return
	}
	auditorium, err := ah.getAuditoriumFromPath(r)
	if err != nil {
		writeAuditoriumError(w, err)
		return
	}
	if auditoriumUpdateReq.CantCols != nil {
		auditorium.CantCols = *auditoriumUpdateReq.CantCols
	}
	if auditoriumUpdateReq.CantRows != nil {
		auditorium.CantRows = *auditoriumUpdateReq.CantRows
	}
	if auditoriumUpdateReq.Name != nil {
		auditorium.Name = *auditoriumUpdateReq.Name
	}
	if err := ah.auditoriumRepo.UpdateAuditorium(r.Context(), auditorium); err != nil {
		if errors.Is(err, database.ErrAuditoriumNotFound) {
			http.Error(w, "auditorium not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error while updating auditorium", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, auditorium, http.StatusOK)
}

func (ah *AuditoriumHandler) getAuditoriumFromPath(r *http.Request) (*domain.Auditorium, error) {
	id := r.PathValue("auditorium_id")
	auditoriumID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	auditorium, err := ah.auditoriumRepo.GetAuditoriumByID(r.Context(), auditoriumID)
	if err != nil {
		return nil, err
	}
	return auditorium, nil
}

func writeAuditoriumError(w http.ResponseWriter, err error) {
	if errors.Is(err, database.ErrAuditoriumNotFound) {
		http.Error(w, "Auditorium not found", http.StatusNotFound)
		return
	}
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		http.Error(w, "invalid auditorium id", http.StatusBadRequest)
		return
	}
	http.Error(w, "error finding auditorium", http.StatusInternalServerError)
}

func (ah *AuditoriumHandler) GetAuditoriums(w http.ResponseWriter, r *http.Request) {
	auditoriums, err := ah.auditoriumRepo.GetAllAuditoriums(r.Context())
	if err != nil {
		http.Error(w, "Error getting auditoriums", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, auditoriums, http.StatusOK)
}

func (ah *AuditoriumHandler) DeleteAuditorium(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("auditorium_id")
	auditoriumID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Wrong id format", http.StatusBadRequest)
		return
	}
	if err := ah.auditoriumRepo.RemoveAuditoriumByID(r.Context(), auditoriumID); err != nil {
		if errors.Is(err, database.ErrAuditoriumNotFound) {
			http.Error(w, "auditorium not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error while removing auditorium", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
