package adminUsecase_test

import (
	"errors"
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/src/admin/adminUsecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAdminRepo struct{}

func (m *mockAdminRepo) SoftDeleteUser(userID string) error {
	if userID == "error" {
		return errors.New("failed to soft delete user")
	}
	return nil
}

func (m *mockAdminRepo) UpdateUser(user adminDto.User) error {
	if user.ID == "error" {
		return errors.New("failed to update user")
	}
	return nil
}

func (m *mockAdminRepo) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	if params.ID == "error" {
		return nil, errors.New("failed to get users by params")
	}
	return []adminDto.User{}, nil
}

func (m *mockAdminRepo) GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	if params.ID == "error" {
		return nil, errors.New("failed to get payment methods by params")
	}
	return []adminDto.PaymentMethod{}, nil
}

func (m *mockAdminRepo) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	if params.ID == "error" {
		return nil, errors.New("failed to get wallets by params")
	}
	return []adminDto.Wallet{}, nil
}

func (m *mockAdminRepo) SavePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	if paymentMethod.PaymentName == "error" {
		return errors.New("failed to save payment method")
	}
	return nil
}

func (m *mockAdminRepo) SoftDeletePaymentMethod(paymentMethodID string) error {
	if paymentMethodID == "error" {
		return errors.New("failed to soft delete payment method")
	}
	return nil
}

func (m *mockAdminRepo) UpdatePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	if paymentMethod.ID == "error" {
		return errors.New("failed to update payment method")
	}
	return nil
}

func (m *mockAdminRepo) GetTransactionRepo(params adminDto.GetTransactionParams) ([]adminDto.GetTransactionResponse, int, error) {
	return []adminDto.GetTransactionResponse{}, 0, nil
}

func TestSoftDeleteUser_Success(t *testing.T) {
	adminRepo := &mockAdminRepo{}
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepo)

	err := adminUsecase.SoftDeleteUser("user123")
	assert.NoError(t, err)
}

func TestSoftDeleteUser_Failure(t *testing.T) {
	adminRepo := &mockAdminRepo{}
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepo)

	err := adminUsecase.SoftDeleteUser("error")
	assert.Error(t, err)
}

func TestUpdateUser_Success(t *testing.T) {
	adminRepo := &mockAdminRepo{}
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepo)

	user := adminDto.UserUpdateRequest{
		ID:          "user123",
		Fullname:    "John Doe",
		Username:    "johndoe",
		Email:       "johndoe@example.com",
		Pin:         "123456",
		PhoneNumber: "123456789",
	}

	err := adminUsecase.UpdateUser(user)
	assert.NoError(t, err)
}

func TestUpdateUser_Failure(t *testing.T) {
	adminRepo := &mockAdminRepo{}
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepo)

	user := adminDto.UserUpdateRequest{
		ID:          "error",
		Fullname:    "John Doe",
		Username:    "johndoe",
		Email:       "johndoe@example.com",
		Pin:         "123456",
		PhoneNumber: "123456789",
	}

	err := adminUsecase.UpdateUser(user)
	assert.Error(t, err)
}
