package adminUsecase

import (
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/pkg/hashingPassword"
	"final-project-enigma/src/admin"
)

type adminUC struct {
	adminRepo admin.AdminRepository
}

func NewAdminUsecase(adminRepo admin.AdminRepository) admin.AdminUsecase {
	return &adminUC{adminRepo}
}

func (u *adminUC) SoftDeleteUser(userID string) error {
	err := u.adminRepo.SoftDeleteUser(userID)
	if err != nil {
		return err
	}
	return nil
}
func (u *adminUC) UpdateUser(request adminDto.UserUpdateRequest) error {
	// Hash the PIN before updating the user
	hashedPin, err := hashingPassword.HashPassword(request.Pin)
	if err != nil {
		return err
	}

	// Create a new user DTO with the hashed PIN
	user := adminDto.User{
		ID:          request.ID,
		Fullname:    request.Fullname,
		Username:    request.Username,
		Email:       request.Email,
		Pin:         hashedPin,
		PhoneNumber: request.PhoneNumber,
	}

	// Update the user with the hashed PIN
	if err := u.adminRepo.UpdateUser(user); err != nil {
		return err
	}
	return nil
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
func (u *adminUC) SavePaymentMethod(request adminDto.CreatePaymentMethod) error {
	paymenMethod := adminDto.PaymentMethod{
		PaymentName: request.PaymentName,
	}

	if err := u.adminRepo.SavePaymentMethod(paymenMethod); err != nil {
		return err
	}
	return nil
}

func (u *adminUC) SoftDeletePaymentMethod(paymenmethodID string) error {
	err := u.adminRepo.SoftDeletePaymentMethod(paymenmethodID)
	if err != nil {
		return err
	}
	return nil
}
func (u *adminUC) UpdatePaymentMethod(request adminDto.UpdatePaymentRequest) error {
	UpdatePaymentMethod := adminDto.PaymentMethod{
		ID:          request.ID,
		PaymentName: request.PaymentName,
	}

	if err := u.adminRepo.UpdatePaymentMethod(UpdatePaymentMethod); err != nil {
		return err
	}

	return nil
}
