package paymentRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/src/payment"
	"fmt"
	"strconv"
	"time"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) payment.PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (repo *paymentRepository) UpdateTransactionStatus(orderID string, status string) error {

	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := repo.db.Exec(query, status, orderID)
	if err != nil {
		return errors.New("failed to update transaction status")
	}

	return err
}

func (repo *paymentRepository) UpdateBalance(orderID, amountStr string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}

	var userID string
	query := `SELECT user_id FROM transactions WHERE id = $1`
	err = repo.db.QueryRow(query, orderID).Scan(&userID)
	if err != nil {
		return err
	}

	var walletID string
	query = `SELECT id FROM wallets WHERE user_id = $1`
	err = repo.db.QueryRow(query, userID).Scan(&walletID)
	if err != nil {
		return err
	}

	query = `UPDATE wallets SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	_, err = repo.db.Exec(query, amount, time.Now(), walletID)
	if err != nil {
		return err
	}

	return nil
}
