package userDelivery

import (
	"final-project-enigma/model/dto/json"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/src/user"

	"github.com/gin-gonic/gin"
)

type userDelivery struct {
	userUC user.UserUsecase
}

func NewUserDelivery(v1Group *gin.RouterGroup, userUC user.UserUsecase) {
	handler := userDelivery{
		userUC: userUC,
	}

	userGroup := v1Group.Group("/user")
	{
		userGroup.GET("/info", middleware.JwtAuthWithRoles("USER"), handler.getDataUser)
		userGroup.POST("/info/upload-image", middleware.JwtAuthWithRoles("USER"), handler.uploadProfilImage)
		userGroup.GET("/info/transactions", middleware.JwtAuthWithRoles("USER"), handler.getTransactionsDetail)
		userGroup.GET("/balance", middleware.JwtAuthWithRoles("USER"), handler.getBalanceInfo)
		userGroup.POST("/balance/topup", middleware.JwtAuthWithRoles("USER"), handler.topupTransactionRequest)
		userGroup.POST("/balance/transfer", middleware.JwtAuthWithRoles("USER"), handler.walletTransactionRequest)
		userGroup.POST("/balance/merchant-payment", middleware.JwtAuthWithRoles("USER"), handler.merchantTransactionRequest)
		userGroup.PUT("/info/update", middleware.JwtAuthWithRoles("USER"), handler.updateDataUser)
		userGroup.DELETE("/delete", middleware.JwtAuthWithRoles("USER"), handler.deletedUser)
	}
}

func (u *userDelivery) updateDataUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	var req userDto.UserUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	err := u.userUC.EditDataUserUC(authHeader, req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, nil, "Success update data", "01", "01")
}

func (u *userDelivery) getDataUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	resp, err := u.userUC.GetDataUserUC(authHeader)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "02", "02")
		return
	}

	json.NewResponSucces(ctx, resp, "Success get data", "01", "01")
}

func (u *userDelivery) uploadProfilImage(ctx *gin.Context) {
	var req userDto.UploadImagesRequest
	authHeader := ctx.GetHeader("Authorization")
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		json.NewResponseError(ctx, "failed to get file", "02", "02")
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		json.NewResponseError(ctx, "failed to open file", "02", "02")
		return
	}
	req.File = file
	err = u.userUC.UploadImagesRequestUC(authHeader, req)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "02", "02")
		return
	}
	json.NewResponSucces(ctx, nil, "Success upload image", "01", "01")

}

func (u *userDelivery) getBalanceInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	resp, err := u.userUC.GetBalanceInfoUC(authHeader)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "02", "02")
		return
	}

	json.NewResponSucces(ctx, resp, "Success get balance info", "01", "01")
}

func (u *userDelivery) getTransactionsDetail(ctx *gin.Context) {
	var params userDto.GetTransactionParams

	authHeader := ctx.GetHeader("Authorization")
	params.TrxId = ctx.Query("trxId")
	params.TrxType = ctx.Query("trxType")
	params.TrxDateStart = ctx.Query("trxDateStart")
	params.TrxDateEnd = ctx.Query("trxDateEnd")
	params.TrxStatus = ctx.Query("paymentStatus")
	params.Page = ctx.Query("page")
	params.Limit = ctx.Query("size")

	resp, totalData, err := u.userUC.GetTransactionUC(authHeader, params)
	if err != nil {
		json.NewResponseForbidden(ctx, "No transaction record", "02", "02")
		return
	}

	json.NewResponSuccesPaging(ctx, resp, "Succes get transaction history", "01", "01", params.Page, totalData)

}

func (u *userDelivery) topupTransactionRequest(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	var req userDto.TopUpTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	resp, err := u.userUC.TopUpTransaction(req, authHeader)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, resp, "create transaction succes", "01", "01")
}

func (u *userDelivery) walletTransactionRequest(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	var req userDto.WalletTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	resp, err := u.userUC.WalletTransaction(req, authHeader)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, resp, "Transfer succes", "01", "01")
}

func (u *userDelivery) merchantTransactionRequest(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	var req userDto.MerchantTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	transactionId, err := u.userUC.MerchantTransaction(req, authHeader)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, transactionId, "Payment merchant success", "01", "01")
}

func (u *userDelivery) deletedUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	err := u.userUC.DeleteUser(authHeader)
	if err != nil {
		if err.Error() == "user not found" {
			json.NewResponseError(ctx, "User not found", "404", "02")
			return
		}
		if err.Error() == "user already deleted" {
			json.NewResponseForbidden(ctx, "User already deleted", "409", "02")
			return
		}
		json.NewResponseError(ctx, "Failed to delete user", "500", "02")
		return
	}

	json.NewResponSucces(ctx, nil, "Success delete your account", "200", "01")
}
