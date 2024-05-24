package transactionDelivery

import (
	"final-project-enigma/model/dto/json"
	"final-project-enigma/model/dto/transactionDtos"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/src/transaction"
	"strconv"

	"github.com/gin-gonic/gin"
)

type transactionDelivery struct {
	transactionUC transaction.TransactionUseCase
}

func NewTransactionDelivery(v1group *gin.RouterGroup, transactionUC transaction.TransactionRepository) {
	handler := transactionDelivery{
		transactionUC: transactionUC,
	}

	transactionGroup := v1group.Group("/transaction")
	{
		transactionGroup.GET("", handler.GetTransaction)
		transactionGroup.GET("/wallet", handler.GetWalletTransaction)
		transactionGroup.GET("/topup", handler.GetTopUpTransaction)
	}
}

func (c *transactionDelivery) GetTransaction(ctx *gin.Context) {
	var req transactionDtos.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
		json.NewResponseError(ctx, "data not found", "01", "02")
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}

	transactions, total, err := c.transactionUC.GetTransaction(page, limit)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "02")
		return
	}

	nextPage := page
	if len(transactions) == limit {
		nextPage++
	}

	paging := gin.H{
		"page":       page,
		"total data": total,
		"next page":  nextPage,
	}

	responseData := gin.H{
		"transactions": transactions,
		"paging":       paging,
	}

	json.NewResponSucces(ctx, responseData, "success", "01", "01")
}

func (c *transactionDelivery) GetWalletTransaction(ctx *gin.Context) {
	var req transactionDtos.WalletTransaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
		json.NewResponseError(ctx, "data not found", "01", "02")
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}

	transactions, total, err := c.transactionUC.GetWalletTransaction(page, limit)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "02")
		return
	}

	nextPage := page
	if len(transactions) == limit {
		nextPage++
	}

	paging := gin.H{
		"page":       page,
		"total data": total,
		"next page":  nextPage,
	}

	ResponseData := gin.H{
		"transactions": transactions,
		"paging":       paging,
	}

	json.NewResponSucces(ctx, ResponseData, "success", "01", "01")
}

func (c *transactionDelivery) GetTopUpTransaction(ctx *gin.Context) {
	var req transactionDtos.TopUpTransaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
		// json.NewResponseError(ctx, "data not found", "01", "01")
		// return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}

	transactions, total, err := c.transactionUC.GetTopUpTransaction(page, limit)
	if err != nil {
		json.NewResponseError(ctx, "data not found", "01", "02")
		return
	}

	nextPage := page
	if len(transactions) == limit {
		nextPage++
	}

	paging := gin.H{
		"page":       nextPage,
		"total data": total,
		"next page":  nextPage,
	}

	ResponseData := gin.H{
		"transactions": transactions,
		"paging":       paging,
	}

	json.NewResponSucces(ctx, ResponseData, "success", "01", "01")
}
