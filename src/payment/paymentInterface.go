package payment

import "final-project-enigma/model/dto/userDto"

type PaymentRepository interface {
	UpdateTransactionStatus(orderID string, status string) error
	UpdateBalance(orderID, amount string) error
}

type PaymentUsecase interface {
	MidtransStatusReq(notification userDto.MidtransNotification) error
}
