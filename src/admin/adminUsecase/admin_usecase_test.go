package adminUsecase_test

import (
	"errors"
	"final-project-enigma/model/dto/adminDto"
	"final-project-enigma/src/admin/adminUsecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) SoftDeleteUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAdminRepository) UpdateUser(user adminDto.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAdminRepository) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	args := m.Called(params)
	return args.Get(0).([]adminDto.User), args.Error(1)
}

func (m *MockAdminRepository) GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	args := m.Called(params)
	return args.Get(0).([]adminDto.PaymentMethod), args.Error(1)
}

func (m *MockAdminRepository) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	args := m.Called(params)
	return args.Get(0).([]adminDto.Wallet), args.Error(1)
}

func (m *MockAdminRepository) SavePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	args := m.Called(paymentMethod)
	return args.Error(0)
}

func (m *MockAdminRepository) SoftDeletePaymentMethod(paymentMethodID string) error {
	args := m.Called(paymentMethodID)
	return args.Error(0)
}

func (m *MockAdminRepository) UpdatePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	args := m.Called(paymentMethod)
	return args.Error(0)
}

func TestGetpaymentMethodByParams_Success(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	expectedPaymentMethods := []adminDto.PaymentMethod{
		{ID: "1", PaymentName: "BCA"},
		{ID: "2", PaymentName: "Gopay"},
	}
	mockRepo.On("GetpaymentMethodByParams", mock.Anything).Return(expectedPaymentMethods, nil)

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	params := adminDto.GetPaymentMethodParams{}
	paymentMethods, err := adminUC.GetpaymentMethodByParams(params)

	assert.NoError(t, err)
	assert.Equal(t, expectedPaymentMethods, paymentMethods)
	mockRepo.AssertExpectations(t)
}

func TestGetpaymentMethodByParams_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("GetpaymentMethodByParams", mock.Anything).Return([]adminDto.PaymentMethod{}, errors.New("some error"))
	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	params := adminDto.GetPaymentMethodParams{}
	paymentMethods, err := adminUC.GetpaymentMethodByParams(params)

	assert.Error(t, err)
	assert.Nil(t, paymentMethods)
	assert.Equal(t, "some error", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetWalletByParams_Success(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	expectedWallets := []adminDto.Wallet{
		{ID: "1", User_id: "1", Balance: "100"},
		{ID: "2", User_id: "2", Balance: "200"},
	}
	mockRepo.On("GetWalletByParams", mock.Anything).Return(expectedWallets, nil)

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	params := adminDto.GetWalletParams{}
	wallets, err := adminUC.GetWalletByParams(params)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallets, wallets)
	mockRepo.AssertExpectations(t)
}

func TestGetWalletByParams_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("GetWalletByParams", mock.Anything).Return([]adminDto.Wallet{}, errors.New("some error"))
	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	params := adminDto.GetWalletParams{}
	wallets, err := adminUC.GetWalletByParams(params)
	assert.Error(t, err)
	assert.Nil(t, wallets)
	assert.Equal(t, "some error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestSoftDeletePaymentMethod_Success(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("SoftDeletePaymentMethod", "test-id").Return(nil)

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	err := adminUC.SoftDeletePaymentMethod("test-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeletePaymentMethod_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("SoftDeletePaymentMethod", "test-id").Return(errors.New("some error"))

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	err := adminUC.SoftDeletePaymentMethod("test-id")

	assert.Error(t, err)
	assert.Equal(t, "some error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdatePaymentMethod_Success(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("UpdatePaymentMethod", mock.Anything).Return(nil)

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	request := adminDto.UpdatePaymentRequest{
		ID:          "test-id",
		PaymentName: "Test Payment",
	}
	err := adminUC.UpdatePaymentMethod(request)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "UpdatePaymentMethod", adminDto.PaymentMethod{
		ID:          "test-id",
		PaymentName: "Test Payment",
	})
	mockRepo.AssertExpectations(t)
}

func TestUpdatePaymentMethod_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("UpdatePaymentMethod", mock.Anything).Return(errors.New("update error"))

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	request := adminDto.UpdatePaymentRequest{
		ID:          "test-id",
		PaymentName: "Test Payment",
	}
	err := adminUC.UpdatePaymentMethod(request)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByParams_Success(t *testing.T) {
	expectedUsers := []adminDto.User{
		{ID: "1", Fullname: "User 1"},
		{ID: "2", Fullname: "User 2"},
	}
	mockRepo := new(MockAdminRepository)
	mockRepo.On("GetUsersByParams", mock.Anything).Return(expectedUsers, nil)

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	params := adminDto.GetUserParams{}
	users, err := adminUC.GetUsersByParams(params)

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByParams_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("GetUsersByParams", mock.Anything).Return([]adminDto.User{}, errors.New("database error"))

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	params := adminDto.GetUserParams{}
	users, err := adminUC.GetUsersByParams(params)

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestSoftDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("SoftDeleteUser", "test-id").Return(nil)

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	err := adminUC.SoftDeleteUser("test-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeleteUser_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("SoftDeleteUser", "test-id").Return(errors.New("some error"))

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	err := adminUC.SoftDeleteUser("test-id")

	assert.Error(t, err)
	assert.Equal(t, "some error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	request := adminDto.UserUpdateRequest{
		ID:          "1",
		Fullname:    "John Doe",
		Username:    "johndoe",
		Email:       "john.doe@example.com",
		Pin:         "123456",
		PhoneNumber: "1234567890",
	}

	mockRepo.On("UpdateUser", mock.Anything).Return(nil)

	err := adminUC.UpdateUser(request)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "UpdateUser", mock.Anything)
}

func TestUpdateUser_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	mockRepo.On("UpdateUser", mock.Anything).Return(errors.New("update error"))

	adminUC := adminUsecase.NewAdminUsecase(mockRepo)
	request := adminDto.UserUpdateRequest{
		ID:          "test-id",
		Fullname:    "Test User",
		Username:    "testuser",
		Email:       "test@example.com",
		Pin:         "1234",
		PhoneNumber: "1234567890",
	}
	err := adminUC.UpdateUser(request)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	mockRepo.AssertExpectations(t)
}
