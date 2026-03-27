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

type MovieHandler struct {
	movieRepo *database.MovieRepository
	cache     *cache.Cache
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

func NewMovieHandler(movieRepo *database.MovieRepository, cache *cache.Cache) *MovieHandler {
	return &MovieHandler{
		cache:     cache,
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
	WriteResponseWithEncoder(w, movie, http.StatusCreated)
}

func (mh *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("movie_id")
	var movie *domain.Movie
	if err := mh.cache.Get(id, &movie); err == nil {
		WriteResponseWithEncoder(w, movie, http.StatusOK)
		return
	}
	movieID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "wrong movie_id format", http.StatusBadRequest)
		return
	}
	movie, err = mh.movieRepo.GetMovieByID(r.Context(), movieID)
	if err != nil {
		if errors.Is(err, database.ErrMovieNotFound) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error getting movie", http.StatusInternalServerError)
		return
	}
	if err := mh.cache.Set(id, movie, time.Hour); err != nil {
		log.Printf("GetMovieByID: failed to cache movie: %v", err)
	}
	WriteResponseWithEncoder(w, movie, http.StatusOK)
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
	WriteResponseWithEncoder(w, movies, http.StatusOK)
}

func parseQueryInt(r *http.Request, key string, defaultValue int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(val)
}

func (mh *MovieHandler) GetMoviesWithProjections(w http.ResponseWriter, r *http.Request) {
	var movies []domain.Movie
	if err := mh.cache.Get("currentMovies", &movies); err != nil {
		movies, err = mh.movieRepo.GetNowShowingMovies(r.Context())
		if err != nil {
			http.Error(w, "error getting movies", http.StatusInternalServerError)
			return
		}
		if err := mh.cache.Set("currentMovies", movies, time.Hour); err != nil {
			log.Printf("GetMoviesWithProjections: failed to update cache: %v", err)
		}
	}

	WriteResponseWithEncoder(w, movies, http.StatusOK)
}

func (mh *MovieHandler) GetFutureMovies(w http.ResponseWriter, r *http.Request) {
	var movies []domain.Movie
	if err := mh.cache.Get("futureMovies", &movies); err != nil {
		movies, err = mh.movieRepo.GetFutureMovies(r.Context())
		if err != nil {
			http.Error(w, "error getting movies", http.StatusInternalServerError)
			return
		}
		if err := mh.cache.Set("futureMovies", movies, time.Hour); err != nil {
			log.Printf("GetFutureMovies: failed to update cache: %v", err)
		}
	}

	WriteResponseWithEncoder(w, movies, http.StatusOK)
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
	var movie *domain.Movie
	if err := mh.cache.Get(id, &movie); err != nil {
		movie, err = mh.movieRepo.GetMovieByID(r.Context(), movieID)
		if err != nil {
			if errors.Is(err, database.ErrMovieNotFound) {
				http.Error(w, "movie not found", http.StatusNotFound)
				return
			}
			http.Error(w, "error finding movie", http.StatusInternalServerError)
			return
		}
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
	if err := mh.cache.Set(id, movie, time.Hour); err != nil {
		log.Printf("UpdateMovie: failed to update cache: %v", err)
	}
	if err := mh.cache.Delete("currentMovies"); err != nil {
		log.Printf("UpdateMovie: failed to invalidate currentMovies cache: %v", err)
	}
	if err := mh.cache.Delete("futureMovies"); err != nil {
		log.Printf("UpdateMovie: failed to invalidate futureMovies cache: %v", err)
	}
	WriteResponseWithEncoder(w, movie, http.StatusOK)
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
	if err := mh.cache.Delete(id); err != nil {
		log.Printf("DeleteMovie: failed to invalidate movie cache: %v", err)
	}
	if err := mh.cache.Delete("currentMovies"); err != nil {
		log.Printf("DeleteMovie: failed to invalidate currentMovies cache: %v", err)
	}
	if err := mh.cache.Delete("futureMovies"); err != nil {
		log.Printf("DeleteMovie: failed to invalidate futureMovies cache: %v", err)
	}
	w.WriteHeader(http.StatusNoContent)
}
