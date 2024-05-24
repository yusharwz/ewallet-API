package authDelivery

import (
	"final-project-enigma/model/dto/json"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/src/auth"
	"fmt"

	"github.com/gin-gonic/gin"
)

type authDelivery struct {
	authUC auth.AuthUsecase
}

func NewAuthDelivery(v1Group *gin.RouterGroup, authUC auth.AuthUsecase) {
	handler := authDelivery{
		authUC: authUC,
	}

	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/register", middleware.BasicAuth, handler.createUserRequest)
		authGroup.POST("/request-otp/email", middleware.BasicAuth, handler.loginUserCodeReuqestEmail)
		authGroup.POST("/request-otp/message", middleware.BasicAuth, handler.loginUserCodeReuqestSMS)
		authGroup.POST("/login", middleware.BasicAuth, handler.loginUserReuqest)
		authGroup.GET("/activate-account", handler.activatedAccount)
		authGroup.POST("/forget-pin", middleware.BasicAuth, handler.forgotPinReq)
		authGroup.POST("/reset-pin", middleware.BasicAuth, handler.resetPin)
	}
}

func (a *authDelivery) loginUserCodeReuqestEmail(ctx *gin.Context) {
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

	err := a.authUC.LoginCodeReqEmail(req.Email)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, nil, "Cek your email inbox", "01", "01")
}

func (a *authDelivery) loginUserCodeReuqestSMS(ctx *gin.Context) {
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

	err := a.authUC.LoginCodeReqSMS(req.PhoneNumber)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, nil, "Cek your message inbox", "01", "01")
}

func (a *authDelivery) loginUserReuqest(ctx *gin.Context) {
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

	token, err := a.authUC.LoginReq(req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, token, "login succes", "01", "01")
}

func (a *authDelivery) createUserRequest(ctx *gin.Context) {
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

	_, err := a.authUC.CreateReq(req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, nil, "create account succes, please check your email for activated your account", "01", "01")
}

func (a *authDelivery) activatedAccount(ctx *gin.Context) {
	var req userDto.ActivatedAccountReq

	req.Email = ctx.Query("email")
	req.Fullname = ctx.Query("fullname")
	req.Unique = ctx.Query("unique")
	req.Code = ctx.Query("code")

	err := a.authUC.ActivatedAccount(req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, nil, "your account has been activated", "01", "01")
}

func (a *authDelivery) forgotPinReq(ctx *gin.Context) {
	var req userDto.FogetPinReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	err := a.authUC.ForgotPinReqUC(req)
	if err != nil {
		fmt.Println(err)
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, nil, "Check your email for reset pin link", "01", "01")
}

func (a *authDelivery) resetPin(ctx *gin.Context) {
	var req userDto.ForgetPinParams

	req.Email = ctx.Query("email")
	req.Username = ctx.Query("username")
	req.Unique = ctx.Query("unique")
	req.Code = ctx.Query("code")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	err := a.authUC.ResetPinUC(req)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponSucces(ctx, nil, "Succes change your pin", "01", "01")
}
