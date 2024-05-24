package transactionUseCase

import (
	"final-project-enigma/model/dto/transactionDtos"
	"final-project-enigma/src/transaction"
	"log"
)

type transactionUC struct {
	transactionRepo transaction.TransactionRepository
}

func (t *transactionUC) GetTransaction(page int, limit int) ([]transactionDtos.Transaction, int, error) {
	results, total, err := t.transactionRepo.GetTransaction(page, limit)
	if err != nil {
		// Log kesalahan secara rinci
		log.Printf("Failed to get transactions from repository: %v", err)
		return nil, 0, err
	}
	return results, total, nil
}

func (t *transactionUC) GetTopUpTransaction(page int, limit int) ([]transactionDtos.TopUpTransaction, int, error) {
	results, total, err := t.transactionRepo.GetTopUpTransaction(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func NewTransactionUseCase(transactionRepo transaction.TransactionRepository) transaction.TransactionUseCase {
	return &transactionUC{transactionRepo}
}

func (t *transactionUC) GetWalletTransaction(page int, limit int) ([]transactionDtos.WalletTransaction, int, error) {
	results, total, err := t.transactionRepo.GetWalletTransaction(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
