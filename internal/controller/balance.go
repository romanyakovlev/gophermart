package controller

import (
	"encoding/json"
	"net/http"

	"github.com/romanyakovlev/gophermart/internal/logger"
	"github.com/romanyakovlev/gophermart/internal/middlewares"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type BalanceController struct {
	withdrawal WithdrawalService
	user       UserService
	logger     *logger.Logger
}

func (b BalanceController) GetBalance(w http.ResponseWriter, r *http.Request) {
	authUser, _ := middlewares.GetUserFromContext(r.Context())
	user, err := b.user.GetUser(authUser.UUID)
	if err != nil {
		b.logger.Debugf("something went wrong: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := models.GetBalanceResponse{Withdrawn: user.Withdrawn, Current: user.Current}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		b.logger.Debugf("cannot encode response JSON body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (b BalanceController) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req models.GetWithdrawRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		b.logger.Debugf("cannot decode request JSON body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	authUser, _ := middlewares.GetUserFromContext(r.Context())
	err := b.withdrawal.CreateWithdrawal(req.Order, req.Sum, authUser.UUID)
	if err != nil {
		b.logger.Debugf("something went wrong: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = b.user.WithdrawUserBalance(req.Sum, authUser.UUID)
	if err != nil {
		b.logger.Debugf("something went wrong: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (b BalanceController) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	authUser, _ := middlewares.GetUserFromContext(r.Context())
	withdrawals, err := b.withdrawal.GetWithdrawalsByUser(authUser.UUID)
	if err != nil {
		b.logger.Debugf("something went wrong: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(withdrawals); err != nil {
		b.logger.Debugf("cannot encode response JSON body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NewBalanceController(withdrawal WithdrawalService, user UserService, logger *logger.Logger) *BalanceController {
	return &BalanceController{withdrawal: withdrawal, user: user, logger: logger}
}
