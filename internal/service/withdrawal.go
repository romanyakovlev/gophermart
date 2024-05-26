package service

import (
	"github.com/google/uuid"

	"github.com/romanyakovlev/gophermart/internal/config"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type WithdrawalService struct {
	config         config.Config
	withdrawalRepo WithdrawalRepository
}

func (w WithdrawalService) CreateWithdrawal(order string, sum float64, userID uuid.UUID) error {
	err := w.withdrawalRepo.Create(order, sum, userID)
	if err != nil {
		return err
	}
	return nil
}

func (w WithdrawalService) GetWithdrawalsByUser(userID uuid.UUID) ([]models.Withdrawal, error) {
	withdrawals, err := w.withdrawalRepo.GetByUser(userID)
	if err != nil {
		return []models.Withdrawal{}, err
	}
	return withdrawals, nil
}

func NewWithdrawalService(config config.Config, withdrawalRepo WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{
		config:         config,
		withdrawalRepo: withdrawalRepo,
	}
}
