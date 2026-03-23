package database

import (
	"cinemasys/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

var ErrUserNotFound = errors.New("User not found")

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, hashed_password, document_number, profile_picture, role) VALUES ($1, $2, $3, $4, $5) returning user_id`
	return ur.db.QueryRowContext(ctx, query, user.Email, user.HashedPassword, user.DocumentNumber, user.ProfilePicture, user.Role).Scan(&user.UserID)
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE user_id = $1`
	return ur.getUserByArg(ctx, query, userID)
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	return ur.getUserByArg(ctx, query, email)
}

func (ur *UserRepository) getUserByArg(ctx context.Context, query string, arg any) (*domain.User, error) {
	var user domain.User
	if err := ur.db.GetContext(ctx, &user, query, arg); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email = $1, hashed_password = $2, document_number = $3, profile_picture = $4, role = $5 WHERE user_id = $6`
	result, err := ur.db.ExecContext(ctx, query, user.Email, user.HashedPassword, user.DocumentNumber, user.ProfilePicture, user.Role, user.UserID)
	return CheckErrResult(result, err, ErrUserNotFound)
}

func (ur *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := ur.db.ExecContext(ctx, query, userID)
	return CheckErrResult(result, err, ErrUserNotFound)
}
