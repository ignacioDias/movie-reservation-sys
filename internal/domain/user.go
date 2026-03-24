package domain

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidEmail error = errors.New("Invalid email")
var ErrInvalidPassword error = errors.New("Invalid password")
var ErrEmptyDNI error = errors.New("Empty DNI")

type User struct {
	Email          string `db:"email" json:"email"`
	HashedPassword string `db:"hashed_password" json:"hashedPassword"`
	DocumentNumber string `db:"document_number" json:"documentNumber"`
	UserID         int64  `db:"user_id" json:"userId"`
	ProfilePicture string `db:"profile_picture" json:"profilePicture"`
	Role           Role   `json:"role" db:"role"`
}

type Role int

const (
	USER Role = iota
	ADMIN
)
const defaultProfilePicture = "/assets/avatars/batman.webp"

func NewUser(email string, password string, documentNumber string, role Role, profilePicture string) (*User, error) {
	if !IsValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	if !IsValidPassword(password) {
		return nil, ErrInvalidPassword
	}
	if !IsValidProfilePicture(profilePicture) {
		profilePicture = defaultProfilePicture
	}
	if documentNumber == "" {
		return nil, ErrEmptyDNI
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:          email,
		HashedPassword: string(hashedPassword),
		ProfilePicture: profilePicture,
		DocumentNumber: documentNumber,
		Role:           role,
	}, nil
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) >= 32 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func IsValidProfilePicture(pp string) bool {
	validAvatars := map[string]bool{
		"/assets/avatars/batman.webp":    true,
		"/assets/avatars/joker.webp":     true,
		"/assets/avatars/spiderman.webp": true,
		"/assets/avatars/dune.webp":      true,
		"/assets/avatars/deniro.webp":    true,
		"/assets/avatars/dicaprio.webp":  true,
		"/assets/avatars/maverick.webp":  true,
		"/assets/avatars/samuel.webp":    true,
		"/assets/avatars/travolta.webp":  true,
	}
	return validAvatars[pp]
}

func (u *User) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
