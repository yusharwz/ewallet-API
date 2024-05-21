package admin

import "final-project-enigma/model/dto/adminDto"

type AdminRepository interface {
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
}

type AdminUsecase interface {
	GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error)
	GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error)
	GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error)
}
