package transactionRepository

import (
	"database/sql"
	"final-project-enigma/model/dto/transactionDtos"
	"final-project-enigma/src/transaction"
	"fmt"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) transaction.TransactionRepository {
	return &transactionRepository{db}
}

func (t *transactionRepository) GetTransaction(page int, limit int) ([]transactionDtos.Transaction, int, error) {
	offset := (page - 1) * limit
	rows, err := t.db.Query("select id, user_id, transaction_type, amount, description, status, created_at from transactions LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data :%w", err)
	}
	defer rows.Close()

	var transactions []transactionDtos.Transaction
	for rows.Next() {
		var data transactionDtos.Transaction
		err := rows.Scan(&data.Id, &data.UserId, &data.TransactionType, &data.Amount, &data.Description, &data.Status, &data.Created_at)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction data : %w", err)
		}

		var transDetails []transactionDtos.TransactionDetail
		rowDetails, err := t.db.Query("select id, transaction_id, wallet_transaction_id,topup_transaction_id, created_at from transactions_detail", data.Id)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to querytransaction detail data: %w", err)
		}
		defer rowDetails.Close()

		for rowDetails.Next() {
			var transDetail transactionDtos.TransactionDetail
			var walletTransactionID, topupTransactionID sql.NullString
			err := rowDetails.Scan(&transDetail.Id, &transDetail.TransactionId, &walletTransactionID, &topupTransactionID, &transDetail.Created_at)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to scan transaction detail data: %w", err)
			}

			if walletTransactionID.Valid {
				transDetail.WalletTransactionId = walletTransactionID.String
			}
			if topupTransactionID.Valid {
				transDetail.TopUpTransactionId = topupTransactionID.String
			}

			transDetails = append(transDetails, transDetail)
		}

		data.TransactionDetail = transDetails

		transactions = append(transactions, data)
	}

	total := len(transactions)

	return transactions, total, nil
}

func (t *transactionRepository) GetWalletTransaction(page int, limit int) ([]transactionDtos.WalletTransaction, int, error) {
	offset := (page - 1) * limit
	rows, err := t.db.Query("select id, transaction_id, from_wallet_id, to_wallet_id, created_at from wallet_transactions limit $1 offset $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query wallet transaction data : %w", err)
	}
	defer rows.Close()

	var transactionsWallet []transactionDtos.WalletTransaction
	for rows.Next() {
		var data transactionDtos.WalletTransaction
		err := rows.Scan(&data.Id, &data.TransactionId, &data.FromWalletId, &data.ToWalletId, &data.Created_at)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan wallet transanction data:%w", err)
		}
		transaction := transactionDtos.WalletTransaction{
			Id:            data.Id,
			TransactionId: data.TransactionId,
			FromWalletId:  data.FromWalletId,
			ToWalletId:    data.ToWalletId,
			Created_at:    data.Created_at,
		}

		transactionsWallet = append(transactionsWallet, transaction)
	}

	total := len(transactionsWallet)
	return transactionsWallet, total, nil
}

func (t *transactionRepository) GetTopUpTransaction(page int, limit int) ([]transactionDtos.TopUpTransaction, int, error) {
	offset := (page - 1) * limit
	rows, err := t.db.Query("select id, transaction_id, payment_method_id, created_at from topup_transactions limit $1 offset $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query get topup transaction data:%w", err)
	}
	defer rows.Close()

	var transactionsTopUp []transactionDtos.TopUpTransaction
	for rows.Next() {
		var data transactionDtos.TopUpTransaction
		err := rows.Scan(&data.Id, &data.TransactionId, &data.PaymentMethodId, &data.Created_at)
		if err != nil {
			return nil, 0, err
		}
		transaction := transactionDtos.TopUpTransaction{
			Id:              data.Id,
			TransactionId:   data.TransactionId,
			PaymentMethodId: data.PaymentMethodId,
			Created_at:      data.Created_at,
		}

		transactionsTopUp = append(transactionsTopUp, transaction)

	}

	total := len(transactionsTopUp)
	return transactionsTopUp, total, nil
}
