package adminDelivery

import (
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/model/dto/json"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/src/admin"
	"strconv"

	"github.com/gin-gonic/gin"
)

type adminDelivery struct {
	adminUsecase admin.AdminUsecase
}

func (d *adminDelivery) SavePaymentMethod(ctx *gin.Context) {
	var req adminDto.CreatePaymentMethod
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}
	if err := d.adminUsecase.SavePaymentMethod(req); err != nil {
		json.NewResponseError(ctx, err.Error(), "failed to add payment method", "01")
		return
	}
	json.NewResponSucces(ctx, req, "succes", "01", "01")
}
func (d *adminDelivery) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	var updateUser adminDto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	updateUser.ID = userID

	if err := d.adminUsecase.UpdateUser(updateUser); err != nil {
		json.NewResponseError(ctx, err.Error(), "failed to update category", "01")
		return
	}

	json.NewResponSucces(ctx, updateUser, "payment method updated successfully", "01", "05")
}
func (d *adminDelivery) SoftDeletePaymentMethod(ctx *gin.Context) {
	paymentMethodID := ctx.Param("id")
	err := d.adminUsecase.SoftDeletePaymentMethod(paymentMethodID)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "03")
		return
	}
	json.NewResponSucces(ctx, nil, "payment method deleted sukses", "01", "03")
}
func (d *adminDelivery) UpdatePaymentMethod(ctx *gin.Context) {
	paymentMethodID := ctx.Param("id")

	var updatePaymentMethod adminDto.UpdatePaymentRequest
	if err := ctx.ShouldBindJSON(&updatePaymentMethod); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	updatePaymentMethod.ID = paymentMethodID

	if err := d.adminUsecase.UpdatePaymentMethod(updatePaymentMethod); err != nil {
		json.NewResponseError(ctx, err.Error(), "failed to update category", "01")
		return
	}

	json.NewResponSucces(ctx, updatePaymentMethod, "payment method updated successfully", "01", "05")
}
func (d *adminDelivery) SaveUser(ctx *gin.Context) {
	var req adminDto.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}
	if err := d.adminUsecase.SaveUser(req); err != nil {
		json.NewResponseError(ctx, err.Error(), "failed to add payment method", "01")
		return
	}
	json.NewResponSucces(ctx, req, "succes", "01", "01")
}
func (d *adminDelivery) SoftDeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	err := d.adminUsecase.SoftDeleteUser(userID)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "03")
		return
	}
	json.NewResponSucces(ctx, nil, "user  deleted sukses", "01", "03")
}
func (d *adminDelivery) GetUsersByParams(c *gin.Context) {
	params := adminDto.GetUserParams{
		ID:          c.Query("id"),
		Fullname:    c.Query("fullname"),
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phone_number"),
		Page:        c.Query("page"),
		Limit:       c.Query("limit"),
		CreateAt:    c.Query("created_at"),
	}

	users, err := d.adminUsecase.GetUsersByParams(params)
	if err != nil {
		json.NewResponseError(c, "failed to get user data", "02", "02")
		return
	}

	json.NewResponSucces(c, users, "Success get data", "01", "01")
}

func (d *adminDelivery) GetpaymentMethodByParams(c *gin.Context) {
	params := adminDto.GetpaymentMethodParams{
		ID:          c.Query("id"),
		PaymentName: c.Query("payment_name"),
		CreatedAt:   c.Query("created_at"),
		Page:        c.Query("page"),
		Limit:       c.Query("limit"),
	}
	paymentMethods, err := d.adminUsecase.GetpaymentMethodByParams(params)
	if err != nil {
		json.NewResponseError(c, "failed to get payment method data", "02", "02")
		return
	}

	json.NewResponSucces(c, paymentMethods, "Success get payment method data", "01", "01")
}

func (d *adminDelivery) GetWalletByParams(c *gin.Context) {
	params := adminDto.GetWalletParams{
		ID:        c.Query("id"),
		User_id:   c.Query("user_id"),
		CreatedAt: c.Query("created_at"),
		Page:      c.Query("page"),
		Limit:     c.Query("limit"),
	}

	if minBalanceStr := c.Query("min_balance"); minBalanceStr != "" {
		minBalance, err := strconv.ParseFloat(minBalanceStr, 64)
		if err != nil {
			json.NewResponseError(c, "invalid min_balance format", "02", "02")
			return
		}
		params.MinBalance = &minBalance
	}

	if maxBalanceStr := c.Query("max_balance"); maxBalanceStr != "" {
		maxBalance, err := strconv.ParseFloat(maxBalanceStr, 64)
		if err != nil {
			json.NewResponseError(c, "invalid max_balance format", "02", "02")
			return
		}
		params.MaxBalance = &maxBalance
	}

	wallets, err := d.adminUsecase.GetWalletByParams(params)
	if err != nil {
		json.NewResponseError(c, "failed to get wallet data", "02", "02")
		return
	}

	json.NewResponSucces(c, wallets, "Success get wallet data", "01", "01")
}

func NewAdminDelivery(router *gin.RouterGroup, adminUsecase admin.AdminUsecase) {
	handler := adminDelivery{adminUsecase: adminUsecase}

	adminGroup := router.Group("/admin")
	{
		adminGroup.Use(middleware.BasicAuth)

		adminGroup.GET("/users", handler.GetUsersByParams)
		adminGroup.POST("/user", handler.SaveUser)
		adminGroup.DELETE("/user/:id", handler.SoftDeleteUser)
		adminGroup.PUT("/user/:id", handler.UpdateUser)

		adminGroup.GET("/paymentMethod", handler.GetpaymentMethodByParams)
		adminGroup.POST("/paymentMethod", handler.SavePaymentMethod)
		adminGroup.PUT("/paymentMethod/:id", handler.UpdatePaymentMethod)
		adminGroup.DELETE("/paymentMethod/:id", handler.SoftDeletePaymentMethod)

		adminGroup.GET("/wallet", handler.GetWalletByParams)
	}
}
