package apperrors

import "errors"

var (
	ErrOrderAlreadyCreatedByRequestedUser = errors.New("order is already created by the requested user")
	ErrOrderAlreadyCreatedByDifferentUser = errors.New("order is already created by a different user")
	ErrOrderNotFound                      = errors.New("order not found")
	ErrUserNotFound                       = errors.New("user not found")
	ErrUserAlreadyExists                  = errors.New("user already exists")
	ErrInvalidLoginOrPassword             = errors.New("invalid login or password")
)
