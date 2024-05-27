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

func NewAdminDelivery(router *gin.RouterGroup, adminUsecase admin.AdminUsecase) {
	handler := adminDelivery{adminUsecase: adminUsecase}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/users", middleware.JwtAuthWithRoles("ADMIN"), handler.GetUsersByParams)
		adminGroup.DELETE("/user/:id", middleware.JwtAuthWithRoles("ADMIN"), handler.SoftDeleteUser)
		adminGroup.PUT("/user/:id", middleware.JwtAuthWithRoles("ADMIN"), handler.UpdateUser)
		adminGroup.GET("/paymentMethod", middleware.JwtAuthWithRoles("ADMIN"), handler.GetpaymentMethodByParams)
		adminGroup.POST("/paymentMethod", middleware.JwtAuthWithRoles("ADMIN"), handler.SavePaymentMethod)
		adminGroup.PUT("/paymentMethod/:id", middleware.JwtAuthWithRoles("ADMIN"), handler.UpdatePaymentMethod)
		adminGroup.DELETE("/paymentMethod/:id", middleware.JwtAuthWithRoles("ADMIN"), handler.SoftDeletePaymentMethod)
		adminGroup.GET("/wallet", handler.GetWalletByParams)
		//transaction
		adminGroup.GET("/transaction", middleware.JwtAuthWithRoles("ADMIN"), handler.GetTransaction)
	}
}

func (d *adminDelivery) SavePaymentMethod(c *gin.Context) {
	var req adminDto.CreatePaymentMethod
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(c, validationError, "bad request", "01", "02")
			return
		}
	}
	if err := d.adminUsecase.SavePaymentMethod(req); err != nil {
		json.NewResponseError(c, err.Error(), "failed to add payment method", "01")
		return
	}
	json.NewResponSucces(c, req, "succes", "01", "01")
}
func (d *adminDelivery) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updateUser adminDto.UserUpdateRequest
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(c, validationError, "bad request", "01", "02")
			return
		}
	}

	updateUser.ID = userID

	if err := d.adminUsecase.UpdateUser(updateUser); err != nil {
		json.NewResponseError(c, err.Error(), "failed to update category", "01")
		return
	}

	json.NewResponSucces(c, updateUser, "user updated successfully", "01", "05")
}
func (d *adminDelivery) SoftDeletePaymentMethod(c *gin.Context) {
	paymentMethodID := c.Param("id")
	err := d.adminUsecase.SoftDeletePaymentMethod(paymentMethodID)
	if err != nil {
		json.NewResponseError(c, err.Error(), "01", "03")
		return
	}
	json.NewResponSucces(c, nil, "payment method deleted sukses", "01", "03")
}
func (d *adminDelivery) UpdatePaymentMethod(c *gin.Context) {
	paymentMethodID := c.Param("id")

	var updatePaymentMethod adminDto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&updatePaymentMethod); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponBadRequest(c, validationError, "bad request", "01", "02")
			return
		}
	}

	updatePaymentMethod.ID = paymentMethodID

	if err := d.adminUsecase.UpdatePaymentMethod(updatePaymentMethod); err != nil {
		json.NewResponseError(c, err.Error(), "failed to update category", "01")
		return
	}

	json.NewResponSucces(c, updatePaymentMethod, "payment method updated successfully", "01", "05")
}

func (d *adminDelivery) SoftDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err := d.adminUsecase.SoftDeleteUser(userID)
	if err != nil {
		json.NewResponseError(c, err.Error(), "01", "03")
		return
	}
	json.NewResponSucces(c, nil, "user  deleted sukses", "01", "03")
}

func (d *adminDelivery) GetUsersByParams(c *gin.Context) {
	params := adminDto.GetUserParams{
		ID:          c.Query("id"),
		Fullname:    c.Query("fullname"),
		Username:    c.Query("username"),
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phoneNumber"),
		Roles:       c.Query("roles"),
		Status:      c.Query("status"),
		StartDate:   c.Query("startDate"),
		EndDate:     c.Query("endDate"),
		Page:        c.Query("page"),
		Limit:       c.Query("limit"),
	}

	users, err := d.adminUsecase.GetUsersByParams(params)
	if err != nil {
		json.NewResponseError(c, "failed to get user data: "+err.Error(), "02", "02")
		return
	}

	json.NewResponSucces(c, users, "Success get data", "01", "01")
}

func (d *adminDelivery) GetpaymentMethodByParams(c *gin.Context) {
	params := adminDto.GetPaymentMethodParams{
		ID:          c.Query("id"),
		PaymentName: c.Query("payment_name"),
		CreatedAt:   c.Query("created_at"),
		Page:        c.Query("page"),
		Limit:       c.Query("limit"),
	}
	paymentMethods, err := d.adminUsecase.GetpaymentMethodByParams(params)
	if err != nil {
		json.NewResponseError(c, "failed to get payment method data: "+err.Error(), "02", "02")
		return
	}

	json.NewResponSucces(c, paymentMethods, "Success get payment method data", "01", "01")
}

func (d *adminDelivery) GetWalletByParams(c *gin.Context) {
	params := adminDto.GetWalletParams{
		ID:        c.Query("id"),
		User_id:   c.Query("user_id"),
		Fullname:  c.Query("fullname"),
		Username:  c.Query("username"),
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
		json.NewResponseError(c, "failed to get wallet data: "+err.Error(), "02", "02")
		return
	}

	json.NewResponSucces(c, wallets, "Success get wallet data", "01", "01")
}

func (d *adminDelivery) GetTransaction(ctx *gin.Context) {
	var params adminDto.GetTransactionParams
	params.UserId = ctx.Query("userId")
	params.TrxId = ctx.Query("trxId")
	params.TrxType = ctx.Query("trxType")
	params.TrxDateStart = ctx.Query("trxDateStart")
	params.TrxDateEnd = ctx.Query("trxDateEnd")
	params.TrxStatus = ctx.Query("paymentStatus")
	params.Page = ctx.Query("page")
	params.Limit = ctx.Query("size")

	resp, totalData, err := d.adminUsecase.GetTransactionUC(params)
	if err != nil {
		json.NewResponseForbidden(ctx, "No transaction record", "02", "02")
		return
	}

	json.NewResponSuccesPaging(ctx, resp, "Succes get transaction history", "01", "01", params.Page, totalData)
}
