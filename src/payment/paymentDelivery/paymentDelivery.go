package paymentDelivery

import (
	"final-project-enigma/model/dto/json"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/payment"

	"github.com/gin-gonic/gin"
)

type paymentDelivery struct {
	paymentUC payment.PaymentUsecase
}

func NewPaymentDelivery(v1Group *gin.RouterGroup, paymentUC payment.PaymentUsecase) {
	handler := paymentDelivery{
		paymentUC: paymentUC,
	}

	paymentGroup := v1Group.Group("/payment")
	{
		paymentGroup.POST("/status/midtrans", handler.midtransStatusRequest)
		paymentGroup.GET("/status/midtrans", handler.midtransStatusRequestGet)
	}
}

func (u *paymentDelivery) midtransStatusRequest(ctx *gin.Context) {
	var notification userDto.MidtransNotification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		json.NewResponseError(ctx, "", "01", "01")
	}

	err := u.paymentUC.MidtransStatusReq(notification)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, nil, "create transaction succes", "01", "01")
}

func (u *paymentDelivery) midtransStatusRequestGet(ctx *gin.Context) {

	json.NewResponSucces(ctx, "PaymentSucces", "", "01", "01")
}
