package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"
)

type Movie struct {
	MovieID        int64     `db:"movie_id" json:"movieId"`
	Title          string    `db:"title" json:"title"`
	Description    string    `db:"description" json:"description"`
	PosterImageURL string    `db:"poster_image_url" json:"posterImageUrl"`
	TrailerURL     string    `db:"trailer_url" json:"trailerUrl"`
	Genres         []Genre   `db:"genres" json:"genres"`
	ReleaseDate    time.Time `db:"release_date" json:"releaseDate"`
}

type Genre string

const (
	Action      Genre = "Action"
	Adult       Genre = "Adult"
	Adventure   Genre = "Adventure"
	Animation   Genre = "Animation"
	Biography   Genre = "Biography"
	Comedy      Genre = "Comedy"
	Crime       Genre = "Crime"
	Documentary Genre = "Documentary"
	Drama       Genre = "Drama"
	Family      Genre = "Family"
	Fantasy     Genre = "Fantasy"
	FilmNoir    Genre = "Film Noir"
	History     Genre = "History"
	Horror      Genre = "Horror"
	Musical     Genre = "Musical"
	Music       Genre = "Music"
	Mystery     Genre = "Mystery"
	Romance     Genre = "Romance"
	SciFi       Genre = "Sci-Fi"
	Short       Genre = "Short"
	Sport       Genre = "Sport"
	TalkShow    Genre = "Talk-Show"
	Thriller    Genre = "Thriller"
	War         Genre = "War"
	Western     Genre = "Western"
)

var validGenres = map[Genre]struct{}{
	Action:      {},
	Adult:       {},
	Adventure:   {},
	Animation:   {},
	Biography:   {},
	Comedy:      {},
	Crime:       {},
	Documentary: {},
	Drama:       {},
	Family:      {},
	Fantasy:     {},
	FilmNoir:    {},
	History:     {},
	Horror:      {},
	Musical:     {},
	Music:       {},
	Mystery:     {},
	Romance:     {},
	SciFi:       {},
	Short:       {},
	Sport:       {},
	TalkShow:    {},
	Thriller:    {},
	War:         {},
	Western:     {},
}

func NewMovie(title, description, posterImageURL string, genres []Genre, trailer string, releaseDate time.Time) (*Movie, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	if len(genres) == 0 {
		return nil, errors.New("Genres are required")
	}
	if !AreValidGenres(genres) {
		return nil, errors.New("invalid genres")
	}
	if !IsValidURL(posterImageURL) {
		return nil, errors.New("Invalid poster image")
	}
	if !IsValidURL(trailer) {
		return nil, errors.New("Invalid trailer")
	}
	return &Movie{
		Title:          title,
		Description:    description,
		PosterImageURL: posterImageURL,
		Genres:         genres,
		TrailerURL:     trailer,
		ReleaseDate:    releaseDate,
	}, nil
}

func IsValidGenre(g Genre) bool {
	_, ok := validGenres[g]
	return ok
}

func AreValidGenres(genres []Genre) bool {
	for _, genre := range genres {
		if !IsValidGenre(genre) {
			return false
		}
	}
	return true
}

func IsValidURL(u string) bool {
	var urlRegex = regexp.MustCompile(`^(https?://)([a-zA-Z0-9\-]+\.)+[a-zA-Z]{2,}(/[^\s]*)?$`)
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	if parsed.Host == "" {
		return false
	}
	return urlRegex.MatchString(u)
}

func (g *Genre) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		*g = Genre(v)
	case string:
		*g = Genre(v)
	default:
		return fmt.Errorf("cannot scan %T into Genre", src)
	}
	return nil
}

func (g Genre) Value() (driver.Value, error) {
	return string(g), nil
}
