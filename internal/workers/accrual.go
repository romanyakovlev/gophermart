package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/romanyakovlev/gophermart/internal/models"
)

var sleepTime = time.Second * 1

type AccrualRequest struct {
	OrderNumber string
	UserID      uuid.UUID
}

type AccrualService interface {
	FetchOrderDetails(orderNumber string) (*models.OrderResponse, error)
}

type OrdersService interface {
	UpdateOrderStatus(orderNumber string, status string) error
	AccrueOrder(orderNumber string, accrual float64) error
}

type UserService interface {
	AccrueUserBalance(accrual float64, userID uuid.UUID) error
}

type AccrualWorker struct {
	accrual             AccrualService
	orders              OrdersService
	user                UserService
	errorChannel        chan error
	accrualRequestsChan chan AccrualRequest
}

func (w *AccrualWorker) StartAccrualWorker(ctx context.Context, workerID int) {
	for {
		select {
		case req := <-w.accrualRequestsChan:
			w.processAccrualRequest(ctx, req)
		case <-ctx.Done():
			return
		}
	}
}

func (w *AccrualWorker) processAccrualRequest(ctx context.Context, req AccrualRequest) {
	go func() {
		w.orders.UpdateOrderStatus(req.OrderNumber, "PROCESSING")
		for {
			orderDetails, err := w.accrual.FetchOrderDetails(req.OrderNumber)
			/*
				TODO: Создать MockedAccrualService для возвращения данных
				var err error
				orderDetails := models.OrderResponse{
					Order:   req.OrderNumber,
					Status:  "PROCESSED",
					Accrual: rand.Float64() * 1000,
				}
			*/
			if err != nil {
				log.Printf("Error fetching order details: %v", err)
				time.Sleep(sleepTime)
				continue
			}

			if orderDetails.Status == "INVALID" {
				log.Printf("Final status reached for order %s: %s", req.OrderNumber, orderDetails.Status)
				w.orders.UpdateOrderStatus(req.OrderNumber, "INVALID")
				break
			}

			if orderDetails.Status == "PROCESSED" {
				log.Printf("Final status reached for order %s: %s", req.OrderNumber, orderDetails.Status)
				w.orders.UpdateOrderStatus(req.OrderNumber, "PROCESSED")
				w.orders.AccrueOrder(req.OrderNumber, orderDetails.Accrual)
				w.user.AccrueUserBalance(orderDetails.Accrual, req.UserID)
				break
			}

			time.Sleep(sleepTime)
		}
	}()
}

func (w *AccrualWorker) StartErrorListener(ctx context.Context) {
	for {
		select {
		case err := <-w.errorChannel:
			fmt.Printf("Error processing request: %v\n", err)
		case <-ctx.Done():
			fmt.Println("Error listener shutting down due to context cancellation.")
			return
		}
	}
}

func InitAccrualWorker(o OrdersService, u UserService, a AccrualService, accrualRequestsChan chan AccrualRequest) *AccrualWorker {
	return &AccrualWorker{
		accrual:             a,
		orders:              o,
		user:                u,
		errorChannel:        make(chan error, 100),
		accrualRequestsChan: accrualRequestsChan,
	}
}
