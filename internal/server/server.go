package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/romanyakovlev/gophermart/internal/config"
	"github.com/romanyakovlev/gophermart/internal/controller"
	"github.com/romanyakovlev/gophermart/internal/db"
	"github.com/romanyakovlev/gophermart/internal/logger"
	"github.com/romanyakovlev/gophermart/internal/middlewares"
	"github.com/romanyakovlev/gophermart/internal/repository"
	"github.com/romanyakovlev/gophermart/internal/service"
	"github.com/romanyakovlev/gophermart/internal/workers"
)

func Router(
	OrdersController *controller.OrdersController,
	UserController *controller.UserController,
	BalanceController *controller.BalanceController,
	sugar *logger.Logger,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.RequestLoggerMiddleware(sugar))
	r.Use(middlewares.GzipMiddleware)
	r.Use(middlewares.JWTMiddleware)
	r.Post("/api/user/register", UserController.UserRegistration)
	r.Post("/api/user/login", UserController.UserLogin)
	r.Post("/api/user/orders", OrdersController.CreateOrder)
	r.Get("/api/user/orders", OrdersController.GetOrders)
	r.Get("/api/user/balance", BalanceController.GetBalance)
	r.Get("/api/user/withdrawals", BalanceController.GetWithdrawals)
	r.Post("/api/user/balance/withdraw", BalanceController.CreateWithdrawal)

	return r
}

func initOrderRepository(serverConfig config.Config, db *sql.DB, sugar *logger.Logger) (service.OrderRepository, error) {
	return repository.NewDBOrder(db)
}

func initWithdrawalRepository(serverConfig config.Config, db *sql.DB, sugar *logger.Logger) (service.WithdrawalRepository, error) {
	return repository.NewDBWithdrawal(db)
}

func initUserRepository(serverConfig config.Config, db *sql.DB, sugar *logger.Logger) (service.UserRepository, error) {
	return repository.NewDBUser(db)
}

func Run() error {
	sugar := logger.GetLogger()
	serverConfig := config.GetConfig(sugar)

	DB, err := db.InitDB(serverConfig.DatabaseURI, sugar)
	if err != nil {
		sugar.Errorf("Server error: %v", err)
		return err
	}
	defer DB.Close()

	orderrepo, err := initOrderRepository(serverConfig, DB, sugar)
	if err != nil {
		sugar.Errorf("Server error: %v", err)
		return err
	}
	withdrawalrepo, err := initWithdrawalRepository(serverConfig, DB, sugar)
	if err != nil {
		sugar.Errorf("Server error: %v", err)
		return err
	}
	userrepo, err := initUserRepository(serverConfig, DB, sugar)
	if err != nil {
		sugar.Errorf("Server error: %v", err)
		return err
	}

	ordersService := service.NewOrdersService(serverConfig, orderrepo)
	withdrawalService := service.NewWithdrawalService(serverConfig, withdrawalrepo)
	userService := service.NewUserService(serverConfig, userrepo)
	accrualService := service.NewAccrualService(serverConfig.AccrualSystemAddress)
	wp := workers.NewWorkerPool(5, *ordersService, *userService, *accrualService)
	OrdersCtrl := controller.NewOrdersController(ordersService, wp, sugar)
	UserCtrl := controller.NewUserController(userService, sugar)
	BalanceCtrl := controller.NewBalanceController(withdrawalService, userService, sugar)
	router := Router(OrdersCtrl, UserCtrl, BalanceCtrl, sugar)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wp.StartAll(ctx)
	//go worker.StartErrorListener(ctx)
	err = http.ListenAndServe(serverConfig.RunAddress, router)
	if err != nil {
		sugar.Errorf("Server error: %v", err)
	}
	return err
}
