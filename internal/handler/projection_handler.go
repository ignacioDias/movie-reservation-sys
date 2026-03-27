package handler

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type ProjectionHandler struct {
	projectionRepo *database.ProjectionRepository
}

type ProjectionRequest struct {
	AuditoriumID    int64                  `json:"auditoriumId"`
	MovieID         int64                  `json:"movieId"`
	ScreeningFormat domain.ScreeningFormat `json:"screeningFormat"`
	Language        domain.Language        `json:"language"`
	StartsAt        time.Time              `json:"startsAt"`
}

type ProjectionUpdateRequest struct {
	AuditoriumID    *int64                  `json:"auditoriumId"`
	MovieID         *int64                  `json:"movieId"`
	ScreeningFormat *domain.ScreeningFormat `json:"screeningFormat"`
	Language        *domain.Language        `json:"language"`
	StartsAt        *time.Time              `json:"startsAt"`
}

func NewProjectionHandler(repo *database.ProjectionRepository) *ProjectionHandler {
	return &ProjectionHandler{projectionRepo: repo}
}
func (ph *ProjectionHandler) CreateProjection(w http.ResponseWriter, r *http.Request) {
	var projectionReq ProjectionRequest
	if err := json.NewDecoder(r.Body).Decode(&projectionReq); err != nil {
		http.Error(w, "wrong format for create projection request", http.StatusBadRequest)
		return
	}
	projection := domain.NewProjection(projectionReq.AuditoriumID, projectionReq.MovieID, projectionReq.ScreeningFormat, projectionReq.Language, projectionReq.StartsAt)
	if err := ph.projectionRepo.CreateProjection(r.Context(), projection); err != nil {
		http.Error(w, "error creating projection", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, projection, http.StatusCreated)
}

func (ph *ProjectionHandler) GetAllProjectionsPerMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("movie_id")
	movieID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for movie id", http.StatusBadRequest)
		return
	}
	projections, err := ph.projectionRepo.GetProjectionsByMovieID(r.Context(), movieID)
	if err != nil {
		http.Error(w, "error getting projections", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, projections, http.StatusOK)
}

func (ph *ProjectionHandler) GetProjection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("projection_id")
	projectionID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for projection id", http.StatusBadRequest)
		return
	}
	projection, err := ph.projectionRepo.GetProjectionByID(r.Context(), projectionID)
	if err != nil {
		if errors.Is(err, database.ErrProjectionNotFound) {
			http.Error(w, "projection not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error getting projection", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, projection, http.StatusOK)
}

func (ph *ProjectionHandler) DeleteProjection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("projection_id")
	projectionID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for projection id", http.StatusBadRequest)
		return
	}
	if err := ph.projectionRepo.DeleteProjection(r.Context(), projectionID); err != nil {
		if errors.Is(err, database.ErrProjectionNotFound) {
			http.Error(w, "projection not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting projection", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
func (ph *ProjectionHandler) UpdateProjection(w http.ResponseWriter, r *http.Request) {
	var updateReq ProjectionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "invalid format in update projection request", http.StatusBadRequest)
		return
	}
	id := r.PathValue("projection_id")
	projectionID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong format for projection id", http.StatusBadRequest)
		return
	}
	projection, err := ph.projectionRepo.GetProjectionByID(r.Context(), projectionID)
	if err != nil {
		if errors.Is(err, database.ErrProjectionNotFound) {
			http.Error(w, "projection not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error finding projection to update", http.StatusInternalServerError)
		return
	}
	if updateReq.AuditoriumID != nil {
		projection.AuditoriumID = *updateReq.AuditoriumID
	}
	if updateReq.Language != nil {
		projection.Language = *updateReq.Language
	}
	if updateReq.MovieID != nil {
		projection.MovieID = *updateReq.MovieID
	}
	if updateReq.ScreeningFormat != nil {
		projection.ScreeningFormat = *updateReq.ScreeningFormat
	}
	if updateReq.StartsAt != nil {
		projection.StartsAt = *updateReq.StartsAt
	}
	if err := ph.projectionRepo.UpdateProjection(r.Context(), projection); err != nil {
		http.Error(w, "error updating projection", http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, projection, http.StatusOK)
}
