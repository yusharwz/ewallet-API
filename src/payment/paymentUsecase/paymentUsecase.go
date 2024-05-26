package paymentUsecase

import (
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/payment"
)

type paymentUC struct {
	paymentRepo payment.PaymentRepository
}

func NewPaymentUsecase(paymentRepo payment.PaymentRepository) payment.PaymentUsecase {
	return &paymentUC{paymentRepo}
}

func (usecase *paymentUC) MidtransStatusReq(notification userDto.MidtransNotification) error {

	switch notification.TransactionStatus {
	case "capture":
		if notification.FraudStatus == "challenge" {
		} else if notification.FraudStatus == "accept" {
			usecase.paymentRepo.UpdateTransactionStatus(notification.OrderID, "settlement")
		}
	case "settlement":
		if err := usecase.paymentRepo.UpdateTransactionStatus(notification.OrderID, "success"); err != nil {
			return err
		}
		if err := usecase.paymentRepo.UpdateBalance(notification.OrderID, notification.GrossAmount); err != nil {
			return err
		}
	case "deny":
		if err := usecase.paymentRepo.UpdateTransactionStatus(notification.OrderID, "deny"); err != nil {
			return err
		}
	case "cancel", "expire":
		if err := usecase.paymentRepo.UpdateTransactionStatus(notification.OrderID, "cancel"); err != nil {
			return err
		}
	case "pending":
		if err := usecase.paymentRepo.UpdateTransactionStatus(notification.OrderID, "pending"); err != nil {
			return err
		}
	default:
		usecase.paymentRepo.UpdateTransactionStatus(notification.OrderID, notification.TransactionStatus)
	}

	return nil
}
