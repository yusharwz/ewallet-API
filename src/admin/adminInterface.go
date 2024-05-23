package admin

import (
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/model/dto/userDto"
)

type AdminRepository interface {
	SoftDeleteUser(userID string) error
	UpdateUser(userID adminDto.User) error
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
	SavePaymentMethod(paymentMethodoe adminDto.PaymentMethod) error
	SoftDeletePaymentMethod(paymentMethodID string) error
	UpdatePaymentMethod(paymenmethodID adminDto.PaymentMethod) error
	UserWalletCreate(id string) (err error)
	UserCreate(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error)
}

type AdminUsecase interface {
	UpdateUser(request adminDto.UserUpdateRequest) error
	SoftDeleteUser(UserID string) error
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
	SavePaymentMethod(request adminDto.CreatePaymentMethod) error
	SoftDeletePaymentMethod(paymentMethodID string) error
	UpdatePaymentMethod(request adminDto.UpdatePaymentRequest) error
	CreateReq(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error)
}
