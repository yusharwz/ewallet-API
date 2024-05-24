package transactionDtos

import "time"

type (
	Transaction struct {
		Id                string              `json:"id"`
		UserId            string              `json:"user_id"`
		TransactionType   string              `json:"transaction_type"`
		Amount            float64             `json:"amount"`
		Description       string              `json:"description"`
		Status            string              `json:"status"`
		Created_at        time.Time           `json:"created_at"`
		TransactionDetail []TransactionDetail `json:"transactions_detail"`
	}
)
