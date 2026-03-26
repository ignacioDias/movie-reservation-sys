package handler

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MovieHandler struct {
	movieRepo *database.MovieRepository
}

type MovieRequest struct {
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	PosterImageURL string         `json:"posterImageUrl"`
	TrailerURL     string         `json:"trailerUrl"`
	Genres         []domain.Genre `json:"genres"`
	ReleaseDate    time.Time      `json:"releaseDate"`
}
type MovieUpdateRequest struct {
	Title          *string        `json:"title"`
	Description    *string        `json:"description"`
	PosterImageURL *string        `json:"posterImageUrl"`
	TrailerURL     *string        `json:"trailerUrl"`
	Genres         []domain.Genre `json:"genres"`
	ReleaseDate    *time.Time     `json:"releaseDate"`
}

const defaultLimit = 20
const defaultOffset = 0

func NewMovieHandler(movieRepo *database.MovieRepository) *MovieHandler {
	return &MovieHandler{
		movieRepo: movieRepo,
	}
}

func (mh *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movieReq MovieRequest
	if err := json.NewDecoder(r.Body).Decode(&movieReq); err != nil {
		http.Error(w, "wrong format: movie from request body", http.StatusBadRequest)
		return
	}
	movie, err := domain.NewMovie(movieReq.Title, movieReq.Description, movieReq.PosterImageURL, movieReq.Genres, movieReq.TrailerURL, movieReq.ReleaseDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := mh.movieRepo.CreateMovie(r.Context(), movie); err != nil {
		http.Error(w, "error while creating movie", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("CreateMovie: failed to encode response: %v", err)
	}
}

func (mh *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("movie_id")
	movieID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong movie_id format", http.StatusBadRequest)
		return
	}
	movie, err := mh.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		if errors.Is(err, database.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error getting movie", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("GetMovieByID: failed to encode response: %v", err)
	}
}

func (mh *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	limit, err := parseQueryInt(r, "limit", defaultLimit)
	if err != nil {
		http.Error(w, "invalid format for limit ", http.StatusBadRequest)
		return
	}
	offset, err := parseQueryInt(r, "offset", defaultOffset)
	if err != nil {
		http.Error(w, "invalid format for offset ", http.StatusBadRequest)
		return
	}
	movies, err := mh.movieRepo.GetAllMovies(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "error getting movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("GetAllMovies: failed to encode response: %v", err)
	}
}

func parseQueryInt(r *http.Request, key string, defaultValue int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}

func (mh *MovieHandler) GetMoviesWithProjections(w http.ResponseWriter, r *http.Request) {
	movies, err := mh.movieRepo.GetNowShowingMovies(r.Context())
	if err != nil {
		http.Error(w, "error getting movies", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("GetMoviesWithProjections: failed to encode response: %v", err)
	}
}

func (mh *MovieHandler) GetFutureMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := mh.movieRepo.GetFutureMovies(r.Context())
	if err != nil {
		http.Error(w, "error getting movies", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Printf("GetFutureMovies: failed to encode response: %v", err)
	}
}

func (mh *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("movie_id")
	movieID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong movie_id format", http.StatusBadRequest)
		return
	}
	var updateReq MovieUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "wrong format for movie update request", http.StatusBadRequest)
		return
	}
	movie, err := mh.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		if errors.Is(err, database.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error finding movie", http.StatusInternalServerError)
		return
	}
	if updateReq.Description != nil {
		if *updateReq.Description == "" {
			http.Error(w, "description cannot be empty", http.StatusBadRequest)
			return
		}
		movie.Description = *updateReq.Description
	}
	if updateReq.Genres != nil {
		if len(updateReq.Genres) == 0 {
			http.Error(w, "genres are required", http.StatusBadRequest)
			return
		}
		if !domain.AreValidGenres(updateReq.Genres) {
			http.Error(w, "invalid genres", http.StatusBadRequest)
			return
		}
		movie.Genres = updateReq.Genres
	}
	if updateReq.PosterImageURL != nil {
		if !domain.IsValidURL(*updateReq.PosterImageURL) {
			http.Error(w, "invalid poster image", http.StatusBadRequest)
			return
		}
		movie.PosterImageURL = *updateReq.PosterImageURL
	}
	if updateReq.ReleaseDate != nil {
		if updateReq.ReleaseDate.IsZero() {
			http.Error(w, "invalid release date", http.StatusBadRequest)
			return
		}
		movie.ReleaseDate = *updateReq.ReleaseDate
	}
	if updateReq.Title != nil {
		if *updateReq.Title == "" {
			http.Error(w, "title cannot be empty", http.StatusBadRequest)
			return
		}
		movie.Title = *updateReq.Title
	}
	if updateReq.TrailerURL != nil {
		if !domain.IsValidURL(*updateReq.TrailerURL) {
			http.Error(w, "invalid trailer", http.StatusBadRequest)
			return
		}
		movie.TrailerURL = *updateReq.TrailerURL
	}
	if err := mh.movieRepo.UpdateMovie(r.Context(), movie); err != nil {
		if errors.Is(err, database.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating movie", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Printf("UpdateMovie: failed to encode response: %v", err)
	}

}

func (mh *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("movie_id")
	movieID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong movie_id format", http.StatusBadRequest)
		return
	}
	if err := mh.movieRepo.DeleteMovie(r.Context(), movieID); err != nil {
		if errors.Is(err, database.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting movie", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
