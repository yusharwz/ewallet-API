package admin

import (
	"final-project-enigma/model/dto/adminDto"
)

type AdminRepository interface {
	UpdateUser(userID adminDto.User) error
	SoftDeleteUser(userID string) error
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
	SavePaymentMethod(paymentMethodoe adminDto.PaymentMethod) error
	SoftDeletePaymentMethod(paymentMethodID string) error
	UpdatePaymentMethod(paymenmethodID adminDto.PaymentMethod) error
	GetTransaction(page int, limit int) ([]adminDto.Transaction, int, error)
	GetWalletTransaction(page int, limit int) ([]adminDto.WalletTransaction, int, error)
	GetTopUpTransaction(page int, limit int) ([]adminDto.TopUpTransaction, int, error)
}

type AdminUsecase interface {
	UpdateUser(request adminDto.UserUpdateRequest) error
	SoftDeleteUser(UserID string) error
	GetUsersByParams(request adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
	SavePaymentMethod(request adminDto.CreatePaymentMethod) error
	SoftDeletePaymentMethod(paymentMethodID string) error
	UpdatePaymentMethod(request adminDto.UpdatePaymentRequest) error
	GetTransaction(page int, limit int) ([]adminDto.Transaction, int, error)
	GetWalletTransaction(page int, limit int) ([]adminDto.WalletTransaction, int, error)
	GetTopUpTransaction(page int, limit int) ([]adminDto.TopUpTransaction, int, error)
}
