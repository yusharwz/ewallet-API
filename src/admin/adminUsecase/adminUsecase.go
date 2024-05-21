package adminUsecase

import (
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/src/admin"
)

type adminUC struct {
	adminRepo admin.AdminRepository
}

func NewAdminUsecase(adminRepo admin.AdminRepository) admin.AdminUsecase {
	return &adminUC{adminRepo}
}

func (u *adminUC) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	users, err := u.adminRepo.GetUsersByParams(params)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (u *adminUC) GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	paymentMethod, err := u.adminRepo.GetpaymentMethodByParams(params)
	if err != nil {
		return nil, err
	}
	return paymentMethod, nil
}
func (u *adminUC) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	wallet, err := u.adminRepo.GetWalletByParams(params)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
