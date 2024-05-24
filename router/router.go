package router

import (
	"database/sql"
	"final-project-enigma/src/transaction/transactionRepository"

	"final-project-enigma/src/transaction/transactionDelivery"
	"final-project-enigma/src/transaction/transactionUseCase"

	"github.com/gin-gonic/gin"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB) {
	transactionRepo := transactionRepository.NewTransactionRepository(db)
	transactionUC := transactionUseCase.NewTransactionUseCase(transactionRepo)
	transactionDelivery.NewTransactionDelivery(v1Group, transactionUC)
}
