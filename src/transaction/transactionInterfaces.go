package transaction

import "final-project-enigma/model/dto/transactionDtos"

type TransactionRepository interface {
	GetTransaction(page int, limit int) ([]transactionDtos.Transaction, int, error)
	GetWalletTransaction(page int, limit int) ([]transactionDtos.WalletTransaction, int, error)
	GetTopUpTransaction(page int, limit int) ([]transactionDtos.TopUpTransaction, int, error)
}

type TransactionUseCase interface {
	GetTransaction(page int, limit int) ([]transactionDtos.Transaction, int, error)
	GetWalletTransaction(page int, limit int) ([]transactionDtos.WalletTransaction, int, error)
	GetTopUpTransaction(page int, limit int) ([]transactionDtos.TopUpTransaction, int, error)
}
