package database

import (
	"cinemasys/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type AuditoriumRepository struct {
	db *sqlx.DB
}

var ErrAuditoriumNotFound = errors.New("Auditorium not found")

func NewAuditoriumRepository(db *sqlx.DB) *AuditoriumRepository {
	return &AuditoriumRepository{db: db}
}

func (ar *AuditoriumRepository) CreateAuditorium(ctx context.Context, auditorium *domain.Auditorium) error {
	query := `INSERT INTO auditoriums (name, cant_rows, cant_cols) VALUES ($1, $2, $3) returning auditorium_id`
	return ar.db.QueryRowContext(ctx, query, auditorium.Name, auditorium.CantRows, auditorium.CantCols).Scan(&auditorium.AuditoriumID)
}

func (ar *AuditoriumRepository) GetAuditoriumByID(ctx context.Context, auditoriumID int64) (*domain.Auditorium, error) {
	query := `SELECT * FROM auditoriums WHERE auditorium_id = $1`
	var auditorium domain.Auditorium
	if err := ar.db.GetContext(ctx, &auditorium, query, auditoriumID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuditoriumNotFound
		}
		return nil, err
	}
	return &auditorium, nil
}

func (ar *AuditoriumRepository) GetAllAuditoriums(ctx context.Context) ([]domain.Auditorium, error) {
	query := `SELECT * FROM auditoriums`
	var auditoriums []domain.Auditorium
	if err := ar.db.SelectContext(ctx, &auditoriums, query); err != nil {
		return nil, err
	}
	return auditoriums, nil
}

func (ar *AuditoriumRepository) UpdateAuditorium(ctx context.Context, auditorium *domain.Auditorium) error {
	query := `UPDATE auditoriums SET name = $1, cant_rows = $2, cant_cols = $3 WHERE auditorium_id = $4`
	result, err := ar.db.ExecContext(ctx, query, auditorium.Name, auditorium.CantRows, auditorium.CantCols, auditorium.AuditoriumID)
	return CheckErrResult(result, err, ErrAuditoriumNotFound)
}

func (ar *AuditoriumRepository) RemoveAuditoriumByID(ctx context.Context, auditoriumID int64) error {
	query := `DELETE FROM auditoriums WHERE auditorium_id = $1`
	result, err := ar.db.ExecContext(ctx, query, auditoriumID)
	return CheckErrResult(result, err, ErrAuditoriumNotFound)
}
