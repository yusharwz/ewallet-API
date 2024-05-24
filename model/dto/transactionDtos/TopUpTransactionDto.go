package transactionDtos

import "time"

type (
	TopUpTransaction struct {
		Id              string    `json:"id"`
		TransactionId   string    `json:"transaction_id"`
		PaymentMethodId string    `json:"payment_method_id"`
		Created_at      time.Time `json:"created_at"`
	}
)
