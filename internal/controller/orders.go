package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/romanyakovlev/gophermart/internal/apperrors"
	"github.com/romanyakovlev/gophermart/internal/logger"
	"github.com/romanyakovlev/gophermart/internal/middlewares"
	"github.com/romanyakovlev/gophermart/internal/workers"
)

type OrdersController struct {
	orders OrdersService
	wp     *workers.WorkerPool
	logger *logger.Logger
}

// TODO: отправить в service layer
func CheckLuhn(number string) bool {
	var sum int
	nDigits := len(number)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if i%2 == parity {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}

func (o OrdersController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	bytes, _ := io.ReadAll(r.Body)
	orderID := string(bytes)
	ok := CheckLuhn(orderID)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	authUser, _ := middlewares.GetUserFromContext(r.Context())
	err := o.orders.CreateOrder(orderID, authUser.UUID)
	if err != nil {
		if errors.Is(err, apperrors.ErrOrderAlreadyCreatedByRequestedUser) {
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, apperrors.ErrOrderAlreadyCreatedByDifferentUser) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			o.logger.Debugf("something went wrong: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	go func() {
		req := workers.AccrualRequest{UserID: authUser.UUID, OrderNumber: orderID}
		o.wp.SendAccrualRequest(req)
	}()

	w.WriteHeader(http.StatusAccepted)
}

func (o OrdersController) GetOrders(w http.ResponseWriter, r *http.Request) {
	authUser, _ := middlewares.GetUserFromContext(r.Context())
	orders, err := o.orders.GetOrdersByUser(authUser.UUID)
	if err != nil {
		o.logger.Debugf("something went wrong: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(orders); err != nil {
		o.logger.Debugf("cannot encode response JSON body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NewOrdersController(orders OrdersService, wp *workers.WorkerPool, logger *logger.Logger) *OrdersController {
	return &OrdersController{orders: orders, wp: wp, logger: logger}
}
