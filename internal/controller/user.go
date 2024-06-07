package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/romanyakovlev/gophermart/internal/apperrors"
	"github.com/romanyakovlev/gophermart/internal/jwt"
	"github.com/romanyakovlev/gophermart/internal/logger"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type UserController struct {
	user   UserService
	logger *logger.Logger
}

func (u *UserController) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var req models.GetAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Invalid request format"})
		return
	}

	user, err := u.user.RegisterUser(req.Login, req.Password)
	if err != nil {
		if err == apperrors.ErrUserAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Login already taken"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	tokenString, err := jwt.GenerateJWTToken(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Error generating token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(1 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}

func (u *UserController) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req models.GetAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Invalid request format"})
		return
	}

	user, err := u.user.AuthenticateUser(req.Login, req.Password)

	if err != nil {
		if err == apperrors.ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Invalid credentials"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Internal server error"})
		}
		return
	}

	tokenString, err := jwt.GenerateJWTToken(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Error generating token"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(1 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}

func NewUserController(user UserService, logger *logger.Logger) *UserController {
	return &UserController{user: user, logger: logger}
}
