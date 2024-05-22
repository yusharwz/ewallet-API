package paymentMethodDto

type (
	CreatePaymentMethod struct {
		PaymentName string `json:"payment_name" binding:"required,max=255"`
	}

	UpdatePaymentRequest struct {
		ID          string `json:"id"`
		PaymentName string `json:"payment_name" binding:"required,max=255"`
	}

	PaymentResponse struct {
		ID          string `json:"id"`
		PaymentName string `json:"payment_name" binding:"required,max=255"`
	}
)
