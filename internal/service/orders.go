package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/romanyakovlev/gophermart/internal/apperrors"
	"github.com/romanyakovlev/gophermart/internal/config"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type OrdersService struct {
	config    config.Config
	orderRepo OrderRepository
}

func (o OrdersService) CreateOrder(orderID string, userID uuid.UUID) error {
	order, err := o.orderRepo.Get(orderID)
	if err == nil {
		if order.UserID == userID {
			return apperrors.ErrOrderAlreadyCreatedByRequestedUser
		} else {
			return apperrors.ErrOrderAlreadyCreatedByDifferentUser
		}
	}

	if !errors.Is(err, apperrors.ErrOrderNotFound) {
		return err
	}
	err = o.orderRepo.Create(orderID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (o OrdersService) GetOrder(orderID string) (models.Order, error) {
	order, err := o.orderRepo.Get(orderID)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (o OrdersService) GetOrdersByUser(userID uuid.UUID) ([]models.Order, error) {
	orders, err := o.orderRepo.GetByUser(userID)
	if err != nil {
		return []models.Order{}, err
	}
	return orders, nil
}

func (o OrdersService) UpdateOrderStatus(orderID string, status string) error {
	err := o.orderRepo.UpdateStatus(orderID, status)
	if err != nil {
		return err
	}
	return nil
}

func (o OrdersService) AccrueOrder(orderID string, accrual float64) error {
	err := o.orderRepo.Accrue(orderID, accrual)
	if err != nil {
		return err
	}
	return nil
}

func NewOrdersService(config config.Config, orderRepo OrderRepository) *OrdersService {
	return &OrdersService{
		config:    config,
		orderRepo: orderRepo,
	}
}
