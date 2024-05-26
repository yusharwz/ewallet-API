package transactionRepository

import (
	"final-project-enigma/model/dto/transactionDtos"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	repo := NewTransactionRepository(mockDB)

	expectedTransactions := []transactionDtos.Transaction{
		{
			Id:              "3446866c-446b-4bc0-8bb2-1b1d27c34176",
			UserId:          "a9fb4734-c1cb-4ccd-aef4-ea414b6b0abe",
			TransactionType: "credit",
			Amount:          100000,
			Description:     "Pembayaran untuk layanan X",
			Status:          "success",
			Created_at:      time.Now(),
			TransactionDetail: []transactionDtos.TransactionDetail{
				{Id: "a7b61173-4c87-4a2f-a534-937feedff3ce", TransactionId: "5cbe4d91-844f-4f5a-90d2-13c0c7c40330", WalletTransactionId: "", TopUpTransactionId: "8e05260b-6a89-45c5-b581-15f02ae8f2b4", Created_at: time.Now()},
				{Id: "2", TransactionId: "2", WalletTransactionId: "", TopUpTransactionId: "2", Created_at: time.Now()},
			},
		},
	}

	page := 1
	limit := 10
	offset := (page - 1) * limit

	rows := sqlmock.NewRows([]string{"id", "user_id", "transaction_type", "amount", "description", "status", "created_at"}).
		AddRow("3446866c-446b-4bc0-8bb2-1b1d27c34176", "a9fb4734-c1cb-4ccd-aef4-ea414b6b0abe", "credit", 100000, "Pembayaran untuk layanan X", "success", time.Now())
	mock.ExpectQuery("select id, user_id, transaction_type, amount, description, status, created_at from transactions LIMIT \\$1 OFFSET \\$2").
		WithArgs(limit, offset).
		WillReturnRows(rows)

	mock.ExpectQuery("select id, transaction_id, wallet_transaction_id, topup_transaction_id, created_at from transactions_detail").
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "wallet_transaction_id", "topup_transaction_id", "created_at"}).
			AddRow("a7b61173-4c87-4a2f-a534-937feedff3ce", "5cbe4d91-844f-4f5a-90d2-13c0c7c40330", "", "8e05260b-6a89-45c5-b581-15f02ae8f2b4", time.Now()).
			AddRow("2", "2", "", "2", time.Now()))

	transactions, total, err := repo.GetTransaction(1, 10)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if total != len(expectedTransactions) {
		t.Errorf("unexpected total count. Expected: %d, Got: %d", len(expectedTransactions), total)
		return
	}

	for i, tr := range transactions {
		expected := expectedTransactions[i]
		if tr.Id != expected.Id || tr.UserId != expected.UserId || tr.TransactionType != expected.TransactionType || tr.Amount != expected.Amount ||
			tr.Description != expected.Description || tr.Status != expected.Status || !tr.Created_at.Equal(expected.Created_at) {
			t.Errorf("unexpected transaction data at index %d. Expected: %+v, Got: %+v", i, expected, tr)
		}

		for j, td := range tr.TransactionDetail {
			expectedDetail := expected.TransactionDetail[j]
			if td.Id != expectedDetail.Id || td.TransactionId != expectedDetail.TransactionId || td.WalletTransactionId != expectedDetail.WalletTransactionId ||
				td.TopUpTransactionId != expectedDetail.TopUpTransactionId || !td.Created_at.Equal(expectedDetail.Created_at) {
				t.Errorf("unexpected transaction detail data at index %d. Expected: %+v, Got: %+v", j, expectedDetail, td)
			}
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetWalletTransactionSuccess(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	repo := NewTransactionRepository(mockDB)

	expectedWalletTransactions := []transactionDtos.WalletTransaction{
		{Id: "1", TransactionId: "1", FromWalletId: "wallet1", ToWalletId: "wallet2", Created_at: time.Now()},
		{Id: "2", TransactionId: "2", FromWalletId: "wallet2", ToWalletId: "wallet3", Created_at: time.Now()},
	}

	page := 1
	pageSize := 10
	offset := (page - 1) * pageSize

	rows := sqlmock.NewRows([]string{"id", "transaction_id", "from_wallet_id", "to_wallet_id", "created_at"}).
		AddRow("1", "1", "wallet1", "wallet2", time.Now()).
		AddRow("2", "2", "wallet2", "wallet3", time.Now())

	mock.ExpectQuery("select id, transaction_id, from_wallet_id, to_wallet_id, created_at from wallet_transactions limit \\$1 offset \\$2").
		WithArgs(pageSize, offset).
		WillReturnRows(rows)

	walletTransactions, totalCount, err := repo.GetWalletTransaction(page, pageSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, expectedWalletTransactions, walletTransactions)
	assert.Equal(t, len(expectedWalletTransactions), totalCount)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTopUpTransaction(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	repo := NewTransactionRepository(mockDB)

	now := time.Now()

	expectedTopupTransactions := []transactionDtos.TopUpTransaction{
		{Id: "1", TransactionId: "1", PaymentMethodId: "1", Created_at: now},
		{Id: "2", TransactionId: "2", PaymentMethodId: "1", Created_at: now},
	}

	page := 1
	pageSize := 10
	offset := (page - 1) * pageSize

	rows := sqlmock.NewRows([]string{"id", "transaction_id", "payment_method_id", "created_at"}).
		AddRow("1", "1", "1", now).
		AddRow("2", "2", "1", now)

	mock.ExpectQuery("select id, transaction_id, payment_method_id, created_at from topup_transactions limit \\$1 offset \\$2").
		WithArgs(pageSize, offset).
		WillReturnRows(rows)

	topupTransactions, totalCount, err := repo.GetTopUpTransaction(page, pageSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, expectedTopupTransactions, topupTransactions)
	assert.Equal(t, len(expectedTopupTransactions), totalCount)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}
