package adminUsecase_test

// import (
// 	"final-project-enigma/model/dto/adminDto"
// 	"final-project-enigma/src/admin/adminUsecase"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // Mock repository
// type MockAdminRepository struct {
// 	mock.Mock
// }

// func (m *MockAdminRepository) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
// 	args := m.Called(params)
// 	return args.Get(0).([]adminDto.User), args.Error(1)
// }

// func (m *MockAdminRepository) GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error) {
// 	args := m.Called(params)
// 	return args.Get(0).([]adminDto.PaymentMethod), args.Error(1)
// }

// func (m *MockAdminRepository) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
// 	args := m.Called(params)
// 	return args.Get(0).([]adminDto.Wallet), args.Error(1)
// }

// func TestGetUsersByParams(t *testing.T) {
// 	mockRepo := new(MockAdminRepository)
// 	usecase := adminUsecase.NewAdminUsecase(mockRepo)
// 	params := adminDto.GetUserParams{}
// 	expectedUsers := []adminDto.User{
// 		{ID: "uuid1", Fullname: "User One", Username: "userone", Email: "userone@example.com"},
// 		{ID: "uuid2", Fullname: "User Two", Username: "usertwo", Email: "usertwo@example.com"},
// 	}

// 	mockRepo.On("GetUsersByParams", params).Return(expectedUsers, nil)

// 	users, err := usecase.GetUsersByParams(params)

// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedUsers, users)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetPaymentMethodByParams(t *testing.T) {
// 	mockRepo := new(MockAdminRepository)
// 	usecase := adminUsecase.NewAdminUsecase(mockRepo)
// 	params := adminDto.GetpaymentMethodParams{}
// 	expectedPaymentMethods := []adminDto.PaymentMethod{
// 		{ID: "uuid1", PaymentName: "ovo"},
// 		{ID: "uuid2", PaymentName: "gopay"},
// 	}

// 	mockRepo.On("GetpaymentMethodByParams", params).Return(expectedPaymentMethods, nil)

// 	paymentMethods, err := usecase.GetpaymentMethodByParams(params)

// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedPaymentMethods, paymentMethods)
// 	mockRepo.AssertExpectations(t)
// }

// func TestGetWalletByParams(t *testing.T) {
// 	mockRepo := new(MockAdminRepository)
// 	usecase := adminUsecase.NewAdminUsecase(mockRepo)
// 	params := adminDto.GetWalletParams{}
// 	expectedWallets := []adminDto.Wallet{
// 		{ID: "uuid1", User_id: "uuid1", Balance: "10000"},
// 		{ID: "uuid2", User_id: "uuid2", Balance: "10000"},
// 	}

// 	mockRepo.On("GetWalletByParams", params).Return(expectedWallets, nil)

// 	wallets, err := usecase.GetWalletByParams(params)

// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedWallets, wallets)
// 	mockRepo.AssertExpectations(t)
// }
