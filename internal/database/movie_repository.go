package database

import (
	"cinemasys/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var ErrMovieNotFound = errors.New("Movie not found")

type MovieRepository struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (mr *MovieRepository) CreateMovie(ctx context.Context, movie *domain.Movie) error {
	query := `INSERT INTO movies (title, description, poster_image_url, trailer_url, genres, release_date) VALUES ($1, $2, $3, $4, $5, $6) returning movie_id`
	return mr.db.QueryRowContext(ctx, query, movie.Title, movie.Description, movie.PosterImageURL, movie.TrailerURL, pq.GenericArray{A: movie.Genres}, movie.ReleaseDate).Scan(&movie.MovieID)
}

func (mr *MovieRepository) GetMovieByID(ctx context.Context, movieID int64) (*domain.Movie, error) {
	query := `SELECT movie_id, title, description, poster_image_url, trailer_url, genres, release_date FROM movies WHERE movie_id = $1`
	var movie domain.Movie
	err := mr.db.QueryRowContext(ctx, query, movieID).Scan(
		&movie.MovieID,
		&movie.Title,
		&movie.Description,
		&movie.PosterImageURL,
		&movie.TrailerURL,
		pq.GenericArray{A: &movie.Genres},
		&movie.ReleaseDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}
	return &movie, nil
}

func (mr *MovieRepository) GetFutureMovies(ctx context.Context) ([]domain.Movie, error) {
	query := `SELECT movie_id, title, description, poster_image_url, trailer_url, genres, release_date FROM movies WHERE release_date > NOW()`
	rows, err := mr.db.QueryContext(ctx, query)
	return mr.getMovies(err, rows)
}

func (mr *MovieRepository) GetAllMovies(ctx context.Context, limit, offset int) ([]domain.Movie, error) {
	query := `SELECT movie_id, title, description, poster_image_url, trailer_url, genres, release_date FROM movies LIMIT $1 OFFSET $2`
	rows, err := mr.db.QueryContext(ctx, query, limit, offset)
	return mr.getMovies(err, rows)
}

func (mr *MovieRepository) GetNowShowingMovies(ctx context.Context) ([]domain.Movie, error) {
	query := ` SELECT DISTINCT m.movie_id, m.title, m.description, m.poster_image_url, m.trailer_url, m.genres, m.release_date
		FROM movies m INNER JOIN projections p ON m.movie_id = p.movie_id WHERE p.starts_at > NOW()`
	rows, err := mr.db.QueryContext(ctx, query)
	return mr.getMovies(err, rows)
}

func (mr *MovieRepository) getMovies(err error, rows *sql.Rows) ([]domain.Movie, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []domain.Movie
	for rows.Next() {
		var movie domain.Movie
		if err := rows.Scan(
			&movie.MovieID,
			&movie.Title,
			&movie.Description,
			&movie.PosterImageURL,
			&movie.TrailerURL,
			pq.GenericArray{A: &movie.Genres},
			&movie.ReleaseDate,
		); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

func (mr *MovieRepository) UpdateMovie(ctx context.Context, movie *domain.Movie) error {
	query := `UPDATE movies SET title = $1, description = $2, poster_image_url = $3, trailer_url = $4, genres = $5, release_date = $6 WHERE movie_id = $7`
	result, err := mr.db.ExecContext(ctx, query, movie.Title, movie.Description, movie.PosterImageURL, movie.TrailerURL, pq.GenericArray{A: movie.Genres}, movie.ReleaseDate, movie.MovieID)
	return CheckErrResult(result, err, ErrMovieNotFound)
}

func (mr *MovieRepository) DeleteMovie(ctx context.Context, movieID int64) error {
	query := `DELETE FROM movies WHERE movie_id = $1`
	result, err := mr.db.ExecContext(ctx, query, movieID)
	return CheckErrResult(result, err, ErrMovieNotFound)
}
