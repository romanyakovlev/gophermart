package controller

import (
	"github.com/google/uuid"

	"github.com/romanyakovlev/gophermart/internal/models"
)

type UserService interface {
	GetUser(userID uuid.UUID) (models.User, error)
	WithdrawUserBalance(sum float64, userID uuid.UUID) error
	AccrueUserBalance(accrual float64, userID uuid.UUID) error
	RegisterUser(login string, password string) (models.User, error)
	AuthenticateUser(login string, password string) (models.User, error)
}

type WithdrawalService interface {
	CreateWithdrawal(order string, sum float64, userID uuid.UUID) error
	GetWithdrawalsByUser(userID uuid.UUID) ([]models.Withdrawal, error)
}

type OrdersService interface {
	CreateOrder(orderID string, userID uuid.UUID) error
	GetOrder(orderID string) (models.Order, error)
	GetOrdersByUser(userID uuid.UUID) ([]models.Order, error)
	UpdateOrderStatus(orderID string, status string) error
	AccrueOrder(orderID string, accrual float64) error
}
