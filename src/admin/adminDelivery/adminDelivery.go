package adminDelivery

import (
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/model/dto/json"
	"final-project-enigma/src/admin"
	"strconv"

	"github.com/gin-gonic/gin"
)

type adminDelivery struct {
	adminUsecase admin.AdminUsecase
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
		adminGroup.GET("/users", handler.GetUsersByParams)
		adminGroup.GET("/paymentMethod", handler.GetpaymentMethodByParams)
		adminGroup.GET("/wallet", handler.GetWalletByParams)
	}
}
