package transactionDtos

import "time"

type (
	TransactionDetail struct {
		Id                  string    `json:"id"`
		TransactionId       string    `json:"transaction_id"`
		WalletTransactionId string    `json:"wallet_transaction_id"`
		TopUpTransactionId  string    `json:"topup_transaction_id"`
		Created_at          time.Time `json:"created_at"`
	}

	// TransactionWalletDetail struct {
	// 	Id                  string    `json:"id"`
	// 	TransactionId       string    `json:"transaction_id"`
	// 	WalletTransactionId string    `json:"wallet_transaction_id"`
	// 	Created_at          time.Time `json:"created_at"`
	// }

	// TransactionTopupDetail struct {
	// 	Id                 string    `json:"id"`
	// 	TransactionId      string    `json:"transaction_id"`
	// 	TopupTransactionId string    `json:"topup_transaction_id"`
	// 	Created_at         time.Time `json:"created_at"`
	// }
)
