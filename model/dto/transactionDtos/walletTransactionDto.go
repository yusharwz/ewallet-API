package transactionDtos

import "time"

type (
	WalletTransaction struct {
		Id            string    `json:"id"`
		TransactionId string    `json:"transaction_id"`
		FromWalletId  string    `json:"from_wallet_id"`
		ToWalletId    string    `json:"to_wallet_id"`
		Created_at    time.Time `json:"created_at"`
	}
)
