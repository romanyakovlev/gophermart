package models

import (
	"time"

	golangjwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthUser struct {
	UUID uuid.UUID
}

type User struct {
	UserID    uuid.UUID
	Current   float64
	Withdrawn float64
	Login     string
	Hash      string
}

type Order struct {
	UserID     uuid.UUID `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Withdrawal struct {
	Number      string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Claims struct {
	UUID uuid.UUID
	golangjwt.RegisteredClaims
}

type OrderUpdate struct {
	Status  *string
	Accrual *float64
}

type OrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type GetBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type GetWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type GetAuthRequest struct {
	Login    string
	Password string
}

type ErrorResponse struct {
	Message string `json:"message"`
}
