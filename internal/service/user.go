package service

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/romanyakovlev/gophermart/internal/apperrors"
	"github.com/romanyakovlev/gophermart/internal/config"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type UserService struct {
	config   config.Config
	userRepo UserRepository
}

func (u UserService) GetUser(userID uuid.UUID) (models.User, error) {
	user, err := u.userRepo.Get(userID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u UserService) WithdrawUserBalance(sum float64, userID uuid.UUID) error {
	err := u.userRepo.WithdrawPoints(sum, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) AccrueUserBalance(accrual float64, userID uuid.UUID) error {
	err := u.userRepo.AccruePoints(accrual, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) RegisterUser(login, password string) (models.User, error) {
	existingUser, _ := u.userRepo.FindByLogin(login)
	if existingUser != (models.User{}) {
		return models.User{}, apperrors.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	newUser := models.User{
		UserID:    uuid.New(),
		Current:   0,
		Withdrawn: 0,
		Login:     login,
		Hash:      string(hashedPassword),
	}
	return u.userRepo.Create(newUser)
}

func (u UserService) AuthenticateUser(login, password string) (models.User, error) {
	user, err := u.userRepo.FindByLogin(login)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return models.User{}, apperrors.ErrInvalidLoginOrPassword
	}

	return user, nil
}

func NewUserService(config config.Config, userRepo UserRepository) *UserService {
	return &UserService{
		config:   config,
		userRepo: userRepo,
	}
}
