package service

import (
	"github.com/google/uuid"

	"github.com/romanyakovlev/gophermart/internal/models"
)

type OrderRepository interface {
	Create(orderID string, userID uuid.UUID) error
	Get(orderID string) (models.Order, error)
	GetByUser(userID uuid.UUID) ([]models.Order, error)
	UpdateStatus(orderID string, status string) error
	Accrue(orderID string, accrual float64) error
}

type WithdrawalRepository interface {
	Create(orderID string, sum float64, userID uuid.UUID) error
	GetByUser(userID uuid.UUID) ([]models.Withdrawal, error)
}

type UserRepository interface {
	Get(userID uuid.UUID) (models.User, error)
	Create(user models.User) (models.User, error)
	FindByLogin(login string) (models.User, error)
	WithdrawPoints(sum float64, userID uuid.UUID) error
	AccruePoints(accrual float64, userID uuid.UUID) error
}
