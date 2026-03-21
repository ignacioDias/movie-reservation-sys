package database

import (
	"cinemasys/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type ProjectionRepository struct {
	db *sqlx.DB
}

var ErrProjectionNotFound = errors.New("Projection not found")

func NewProjectionRepository(db *sqlx.DB) *ProjectionRepository {
	return &ProjectionRepository{db: db}
}

func (pr *ProjectionRepository) CreateProjection(ctx context.Context, projection *domain.Projection) error {
	query := `INSERT INTO projections (auditorium_id, movie_id, screening_format, language, starts_at) VALUES ($1, $2, $3, $4, $5) returning projection_id`
	return pr.db.QueryRowContext(ctx, query, projection.AuditoriumID, projection.MovieID, projection.ScreeningFormat, projection.Language, projection.StartsAt).Scan(&projection.ProjectionID)
}

func (pr *ProjectionRepository) GetProjectionByID(ctx context.Context, projectionID int64) (*domain.Projection, error) {
	query := `SELECT * FROM projections WHERE projection_id = $1`
	var projection domain.Projection
	err := pr.db.GetContext(ctx, &projection, query, projectionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectionNotFound
		}
		return nil, err
	}
	return &projection, nil
}

func (pr *ProjectionRepository) GetProjectionsByMovieID(ctx context.Context, movieID int64) ([]domain.Projection, error) {
	query := `SELECT * FROM projections WHERE movie_id = $1 AND starts_at > NOW()`
	var projections []domain.Projection
	if err := pr.db.SelectContext(ctx, &projections, query, movieID); err != nil {
		return nil, err
	}
	return projections, nil
}

func (pr *ProjectionRepository) UpdateProjection(ctx context.Context, projection *domain.Projection) error {
	query := `UPDATE projections SET auditorium_id = $1, movie_id = $2, screening_format = $3, language = $4, starts_at = $5 WHERE projection_id = $6`
	result, err := pr.db.ExecContext(ctx, query, projection.AuditoriumID, projection.MovieID, projection.ScreeningFormat, projection.Language, projection.StartsAt, projection.ProjectionID)
	return CheckErrResult(result, err, ErrProjectionNotFound)
}

func (pr *ProjectionRepository) DeleteProjection(ctx context.Context, projectionID int64) error {
	query := `DELETE FROM projections WHERE projection_id = $1`
	result, err := pr.db.ExecContext(ctx, query, projectionID)
	return CheckErrResult(result, err, ErrProjectionNotFound)
}
