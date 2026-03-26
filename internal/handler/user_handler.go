package handler

import (
	"cinemasys/internal/database"
	"cinemasys/internal/domain"
	"cinemasys/internal/middleware"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

type UserHandler struct {
	userRepo    *database.UserRepository
	sessionRepo *database.SessionRepository
}

type UserRegisterRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	DocumentNumber string `json:"documentNumber"`
	ProfilePicture string `json:"profilePicture"`
}

type UpdateUserRequest struct {
	Email          *string `json:"email"`
	DocumentNumber *string `json:"documentNumber"`
	ProfilePicture *string `json:"profilePicture"`
	Password       *string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var isProduction bool = os.Getenv("ENV") == "production"

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
	user, err := domain.NewUser(userRequest.Email, userRequest.Password, userRequest.DocumentNumber, domain.USER, userRequest.ProfilePicture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uh.userRepo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "wrong input for login request", http.StatusBadRequest)
		return
	}
	user, err := uh.userRepo.GetUserByEmail(r.Context(), loginReq.Email)
	if errors.Is(err, database.ErrUserNotFound) {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}
	if !user.ComparePasswords(loginReq.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	session := domain.NewSession(user.UserID)
	if err := uh.sessionRepo.CreateSession(r.Context(), session); err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction, // false only in localhost dev
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("LoginUser: failed to encode response: %v", err)
	}
}
func (uh *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := uh.sessionRepo.DeleteSessionsByUserID(r.Context(), userID); err != nil {
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
	user, err := uh.getUserFromPath(r)
	if err != nil {
		writeUserError(w, err)
		return
	}
	user.Role = domain.ADMIN
	if err := uh.userRepo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := uh.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("GetCurrentUser: failed to encode response: %v", err)
	}
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "wrong format for update user", http.StatusBadRequest)
		return
	}
	user, err := uh.getUserFromPath(r)
	if err != nil {
		writeUserError(w, err)
		return
	}
	if req.Email != nil {
		if !domain.IsValidEmail(*req.Email) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			return
		}
		user.Email = *req.Email
	}
	if req.DocumentNumber != nil {
		if *req.DocumentNumber == "" {
			http.Error(w, "document number cannot be empty", http.StatusBadRequest)
			return
		}
		user.DocumentNumber = *req.DocumentNumber
	}
	if req.ProfilePicture != nil {
		if !domain.IsValidProfilePicture(*req.ProfilePicture) {
			http.Error(w, "invalid picture", http.StatusBadRequest)
			return
		}
		user.ProfilePicture = *req.ProfilePicture
	}
	if req.Password != nil {
		if !domain.IsValidPassword(*req.Password) {
			http.Error(w, "invalid password", http.StatusBadRequest)
			return
		}
		hashedPassword, err := domain.HashPassword(*req.Password)
		if err != nil {
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}
		user.HashedPassword = string(hashedPassword)
	}
	if err := uh.userRepo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("UpdateUser: failed to encode response: %v", err)
	}

}

func (uh *UserHandler) getUserFromPath(r *http.Request) (*domain.User, error) {
	id := r.PathValue("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	user, err := uh.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func writeUserError(w http.ResponseWriter, err error) {
	if errors.Is(err, database.ErrUserNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	http.Error(w, "error finding user", http.StatusInternalServerError)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if err := uh.userRepo.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error while deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := uh.userRepo.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error while deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
