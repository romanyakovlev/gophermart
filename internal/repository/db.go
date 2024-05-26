package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/romanyakovlev/gophermart/internal/apperrors"
	"github.com/romanyakovlev/gophermart/internal/models"
)

type DBOrder struct {
	db *sql.DB
}

type DBWithdrawal struct {
	db *sql.DB
}
type DBUser struct {
	db *sql.DB
}

func (r DBOrder) Create(orderID string, userID uuid.UUID) error {
	query := "INSERT INTO orders (number, status, accrual, uploaded_at, user_id) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, orderID, "NEW", 0, time.Now(), userID)
	if err != nil {
		return err
	}
	return nil
}

func (r DBOrder) Get(orderID string) (models.Order, error) {
	var order models.Order
	row := r.db.QueryRow("SELECT user_id, number, status, accrual, uploaded_at FROM orders WHERE number = $1", orderID)
	err := row.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Order{}, apperrors.ErrOrderNotFound
		}
		return models.Order{}, err
	}
	return order, nil
}

func (r DBOrder) GetByUser(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	rows, err := r.db.Query("SELECT user_id, number, status, accrual, uploaded_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, nil
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, nil
	}

	return orders, nil
}

func (r DBOrder) UpdateStatus(orderID string, status string) error {
	query := "UPDATE orders SET status = $1 WHERE number = $2"
	result, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}

func (r DBOrder) Accrue(orderID string, accrual float64) error {
	query := "UPDATE orders SET accrual = $1 WHERE number = $2"
	result, err := r.db.Exec(query, accrual, orderID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}

func (r DBWithdrawal) Create(orderID string, sum float64, userID uuid.UUID) error {
	query := "INSERT INTO withdrawals (number, sum, processed_at, user_id) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, orderID, sum, time.Now(), userID)
	if err != nil {
		return err
	}
	return nil
}

func (r DBWithdrawal) GetByUser(userID uuid.UUID) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	rows, err := r.db.Query("SELECT number, sum, processed_at FROM withdrawals WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal models.Withdrawal
		if err := rows.Scan(&withdrawal.Number, &withdrawal.Sum, &withdrawal.ProcessedAt); err != nil {
			return nil, nil
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, nil
	}

	return withdrawals, nil
}

func (r DBUser) Get(userID uuid.UUID) (models.User, error) {
	var user models.User
	row := r.db.QueryRow("SELECT userid, pass_hash, current_balance, withdrawn_balance, userlogin FROM users WHERE userid = $1", userID)
	err := row.Scan(&user.UserID, &user.Hash, &user.Current, &user.Withdrawn, &user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, apperrors.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (r DBUser) FindByLogin(login string) (models.User, error) {
	var user models.User
	row := r.db.QueryRow("SELECT userid, pass_hash, current_balance, withdrawn_balance, userlogin FROM users WHERE userlogin = $1", login)
	err := row.Scan(&user.UserID, &user.Hash, &user.Current, &user.Withdrawn, &user.Login)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r DBUser) WithdrawPoints(sum float64, userID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	var withdrawnBalance float64
	row := tx.QueryRow("SELECT current_balance, withdrawn_balance FROM users WHERE userid = $1 FOR UPDATE", userID)
	err = row.Scan(&currentBalance, &withdrawnBalance)
	if err != nil {
		return err
	}

	if currentBalance < sum {
		return errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE users SET withdrawn_balance = withdrawn_balance + $1 WHERE userid = $2", sum, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET current_balance = current_balance - $1 WHERE userid = $2", sum, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r DBUser) AccruePoints(accrual float64, userID uuid.UUID) error {
	query := "UPDATE users SET current_balance = current_balance + $1 WHERE userid = $2"
	result, err := r.db.Exec(query, accrual, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}
	return nil
}

func (r DBUser) Create(user models.User) (models.User, error) {
	query := "INSERT INTO users (userid, pass_hash, current_balance, withdrawn_balance, userlogin) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, user.UserID, user.Hash, user.Current, user.Withdrawn, user.Login)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func NewDBOrder(db *sql.DB) (*DBOrder, error) {
	return &DBOrder{db: db}, nil
}

func NewDBWithdrawal(db *sql.DB) (*DBWithdrawal, error) {
	return &DBWithdrawal{db: db}, nil
}

func NewDBUser(db *sql.DB) (*DBUser, error) {
	return &DBUser{db: db}, nil
}
