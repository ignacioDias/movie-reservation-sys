package middleware

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"net/http"
)

type AuthorizationMiddleware struct {
	userRepo *database.UserRepository
}

func NewAuthorizationMiddleware(userRepo *database.UserRepository) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		userRepo: userRepo,
	}
}

func (am *AuthorizationMiddleware) AuthorizationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := am.userRepo.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if user.Role != domain.ADMIN {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
