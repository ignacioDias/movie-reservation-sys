package domain

import (
	"time"
)

type Projection struct {
	ProjectionID    int64           `db:"projection_id" json:"projectionId"`
	AuditoriumID    int64           `db:"auditorium_id" json:"auditoriumId"`
	MovieID         int64           `db:"movie_id" json:"movieId"`
	ScreeningFormat ScreeningFormat `db:"screening_format" json:"screeningFormat"`
	Language        Language        `db:"language" json:"language"`
	StartsAt        time.Time       `db:"starts_at" json:"startsAt"`
}

type ScreeningFormat string

const (
	Format2D ScreeningFormat = "2D"
	Format3D ScreeningFormat = "3D"
)

type Language string

const (
	Spanish  Language = "Spanish"
	Original Language = "Original"
	Other    Language = "Other"
)

func NewProjection(auditoriumID, movieID int64, screeningFormat ScreeningFormat, lang Language, startsAt time.Time) *Projection {
	return &Projection{
		AuditoriumID:    auditoriumID,
		MovieID:         movieID,
		ScreeningFormat: screeningFormat,
		Language:        lang,
		StartsAt:        startsAt,
	}
}
