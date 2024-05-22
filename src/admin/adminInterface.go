package admin

import (
	"final-project-enigma/model/dto/adminDto"
)

type AdminRepository interface {
	SaveUser(user adminDto.User) error
	SoftDeleteUser(userID string) error
	UpdateUser(userID adminDto.User) error
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
	SavePaymentMethod(paymentMethodoe adminDto.PaymentMethod) error
	SoftDeletePaymentMethod(paymentMethodID string) error
	UpdatePaymentMethod(paymenmethodID adminDto.PaymentMethod) error
}

type AdminUsecase interface {
	SaveUser(request adminDto.UserCreateRequest) error
	UpdateUser(request adminDto.UserUpdateRequest) error
	SoftDeleteUser(UserID string) error
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
	SavePaymentMethod(request adminDto.CreatePaymentMethod) error
	SoftDeletePaymentMethod(paymentMethodID string) error
	UpdatePaymentMethod(request adminDto.UpdatePaymentRequest) error
}
