package userDelivery

import (
	"final-project-enigma/model/dto/json"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/src/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userDelivery struct {
	userUC user.UserUsecase
}

func NewUserDelivery(v1Group *gin.RouterGroup, userUC user.UserUsecase) {
	handler := userDelivery{
		userUC: userUC,
	}
	userGroup := v1Group.Group("/users")

	{
		userGroup.POST("/register", middleware.BasicAuth, handler.createUserRequest)
		userGroup.POST("/reqcode/email", middleware.BasicAuth, handler.loginUserCodeReuqestEmail)
		userGroup.POST("/reqcode/whatsapp", middleware.BasicAuth, handler.loginUserCodeReuqestSMS)
		userGroup.POST("/login", middleware.BasicAuth, handler.loginUserReuqest)
		userGroup.GET("/info", middleware.JWTAuth(), handler.getDataUser)
		userGroup.GET("/info/balance", middleware.JWTAuth(), handler.getBalanceInfo)
		userGroup.GET("/info/transactions", middleware.JWTAuth(), handler.getTransactionsDetail)
		userGroup.POST("/info/balance/topup", middleware.JWTAuth(), handler.topupTransactionRequest)
	}
}

func (u *userDelivery) loginUserCodeReuqestEmail(ctx *gin.Context) {
	var req userDto.UserLoginCodeRequestEmail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	err := u.userUC.LoginCodeReqEmail(req.Email)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, nil, "Cek your email inbox", "01", "01")
}

func (u *userDelivery) loginUserCodeReuqestSMS(ctx *gin.Context) {
	var req userDto.UserLoginCodeRequestPhoneNumber
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	err := u.userUC.LoginCodeReqSMS(req.PhoneNumber)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, nil, "Cek your message inbox", "01", "01")
}

func (u *userDelivery) loginUserReuqest(ctx *gin.Context) {
	var req userDto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	token, err := u.userUC.LoginReq(req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, token, "login succes", "01", "01")
}

func (u *userDelivery) createUserRequest(ctx *gin.Context) {
	var req userDto.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	resp, err := u.userUC.CreateReq(req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, resp, "login succes", "01", "01")
}

func (u *userDelivery) getDataUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	resp, err := u.userUC.GetDataUserUC(authHeader)
	if err != nil {
		json.NewResponseError(ctx, "failed to get user data", "02", "02")
		return
	}

	json.NewResponSucces(ctx, resp, "Succes get data", "01", "01")
}

func (u *userDelivery) getBalanceInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	resp, err := u.userUC.GetBalanceInfoUC(authHeader)
	if err != nil {
		json.NewResponseError(ctx, "failed to get user data", "02", "02")
		return
	}

	json.NewResponSucces(ctx, resp, "Succes get balance info", "01", "01")
}

func (u *userDelivery) getTransactionsDetail(ctx *gin.Context) {
	var params userDto.GetTransactionParams

	authHeader := ctx.GetHeader("Authorization")
	params.TrxId = ctx.Query("trxId")
	params.TrxType = ctx.Query("trxType")
	params.TrxDateStart = ctx.Query("transactionDateStart")
	params.TrxDateEnd = ctx.Query("transactionDateEnd")
	params.TrxStatus = ctx.Query("paymentStatus")
	params.Page = ctx.Query("page")
	params.Limit = ctx.Query("limit")

	resp, err := u.userUC.GetTransactionUC(authHeader, params)
	if err != nil {
		fmt.Println(err)
		json.NewResponseError(ctx, "failed to get transaction data", "02", "02")
		return
	}

	json.NewResponSucces(ctx, resp, "Succes get transaction data", "01", "01")

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
