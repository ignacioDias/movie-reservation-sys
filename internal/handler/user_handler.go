package handler

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"cinemasys/internal/middleware"
	"encoding/json"
	"net/http"
	"os"
)

type UserHandler struct {
	userRepo    *database.UserRepository
	sessionRepo *database.SessionRepository
}

type UserRegisterRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	DocumentNumber string `db:"document_number" json:"documentNumber"`
	ProfilePicture string `db:"profile_picture" json:"profilePicture"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var isProduction bool = os.Getenv("ENV") == "production"

const defaultProfilePicture = "/assets/avatars/batman.webp"

func NewUserHandler(userRepo *database.UserRepository, sessionRepo *database.SessionRepository) *UserHandler {
	return &UserHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	if !isValidProfilePicture(userRequest.ProfilePicture) {
		userRequest.ProfilePicture = defaultProfilePicture
	}
	user, err := domain.NewUser(userRequest.Email, userRequest.Password, userRequest.DocumentNumber, domain.USER, userRequest.ProfilePicture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := uh.userRepo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func isValidProfilePicture(pp string) bool {
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

func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := uh.userRepo.GetUserByEmail(r.Context(), loginReq.Email)
	if err == database.ErrUserNotFound {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.ComparePasswords(loginReq.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	session := domain.NewSession(user.UserID)
	if err := uh.sessionRepo.CreateSession(r.Context(), session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction, // false only in localhost dev
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func (uh *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if uh.sessionRepo.DeleteSessionsByUserID(r.Context(), userID) != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
}
func (uh *UserHandler) MakeUserAdmin(w http.ResponseWriter, r *http.Request) {

}
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}
