package adminUsecase_test

import (
	"errors"
	"final-project-enigma/model/dto/adminDto"

	// "final-project-enigma/pkg/helper/hashingPassword"
	"final-project-enigma/src/admin/adminUsecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// Mocking the AdminRepository
type MockAdminRepository struct {
	mock.Mock
}

// GetUsersByParams implements admin.AdminRepository.
func (m *MockAdminRepository) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	args := m.Called(params)
	return args.Get(0).([]adminDto.User), args.Error(1)
}

// GetWalletByParams implements admin.AdminRepository.
func (m *MockAdminRepository) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	args := m.Called(params)
	return args.Get(0).([]adminDto.Wallet), args.Error(1)
}

// GetpaymentMethodByParams implements admin.AdminRepository.
func (m *MockAdminRepository) GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	args := m.Called(params)
	return args.Get(0).([]adminDto.PaymentMethod), args.Error(1)
}

// SavePaymentMethod implements admin.AdminRepository.
func (m *MockAdminRepository) SavePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	args := m.Called(paymentMethod)
	return args.Error(0)
}

// SoftDeletePaymentMethod implements admin.AdminRepository.
func (m *MockAdminRepository) SoftDeletePaymentMethod(paymentMethodID string) error {
	args := m.Called(paymentMethodID)
	return args.Error(0)
}

// UpdatePaymentMethod implements admin.AdminRepository.
func (m *MockAdminRepository) UpdatePaymentMethod(paymenmethodID adminDto.PaymentMethod) error {
	args := m.Called(paymenmethodID)
	return args.Error(0)
}

// UpdateUser implements admin.AdminRepository.
func (m *MockAdminRepository) UpdateUser(user adminDto.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAdminRepository) SoftDeleteUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestSoftDeleteUser_Succes(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		userID := "123"
		mockRepo.On("SoftDeleteUser", userID).Return(nil)

		err := usecase.SoftDeleteUser(userID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
func TestSoftDeleteUser_Failure(t *testing.T) {
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	t.Run("failure", func(t *testing.T) {
		userID := "123"
		mockRepo.On("SoftDeleteUser", userID).Return(errors.New("error deleting user"))

		err := usecase.SoftDeleteUser(userID)
		assert.Error(t, err)
		assert.Equal(t, "error deleting user", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
func TestUpdateUser_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	request := adminDto.UserUpdateRequest{
		ID:          "123",
		Fullname:    "John Doe",
		Username:    "johndoe",
		Email:       "john@example.com",
		Pin:         "123456", // Kata sandi belum di-hash untuk pengujian ini
		PhoneNumber: "123456789",
	}

	// Ekspektasi pemanggilan UpdateUser di repository
	mockRepo.On("UpdateUser", mock.Anything).Return(nil)

	// Pengujian
	err := usecase.UpdateUser(request)

	// Verifikasi
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	request := adminDto.UserUpdateRequest{
		ID:          "123",
		Fullname:    "John Doe",
		Username:    "johndoe",
		Email:       "john@example.com",
		Pin:         "123456", // Kata sandi belum di-hash untuk pengujian ini
		PhoneNumber: "123456789",
	}

	// Simulasikan kesalahan dari repository
	expectedErr := errors.New("database error")
	mockRepo.On("UpdateUser", mock.Anything).Return(expectedErr)

	// Pengujian
	err := usecase.UpdateUser(request)

	// Verifikasi
	assert.EqualError(t, err, expectedErr.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByParams_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	params := adminDto.GetUserParams{
		// Isi parameter sesuai kebutuhan pengujian
		Username: "john_doe",         // Filter berdasarkan username "john_doe"
		Email:    "john@example.com", // Filter berdasarkan email "john@example.com"
	}

	expectedUsers := []adminDto.User{
		// Isi dengan data pengguna yang diharapkan
		{
			ID:          "1",
			Fullname:    "John Doe",
			Username:    "john_doe",
			Email:       "john@example.com",
			Pin:         "hashed_password_1",
			PhoneNumber: "123456789",
		},
		{
			ID:          "2",
			Fullname:    "Jane Smith",
			Username:    "jane_smith",
			Email:       "jane@example.com",
			Pin:         "hashed_password_2",
			PhoneNumber: "987654321",
		},
	}

	// Ekspektasi pemanggilan GetUsersByParams di repository
	mockRepo.On("GetUsersByParams", params).Return(expectedUsers, nil)

	// Pengujian
	users, err := usecase.GetUsersByParams(params)

	// Verifikasi
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestGetUsersByParams_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	params := adminDto.GetUserParams{
		Username: "john_doe",         // Filter berdasarkan username "john_doe"
		Email:    "john@example.com", // Filter berdasarkan email "john@example.com"
	}

	// Simulasikan kesalahan dari repository
	expectedErr := errors.New("database error")
	mockRepo.On("GetUsersByParams", params).Return([]adminDto.User{}, expectedErr)

	// Pengujian
	users, err := usecase.GetUsersByParams(params)

	// Verifikasi
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.EqualError(t, err, expectedErr.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetpaymentMethodByParams_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	params := adminDto.GetpaymentMethodParams{
		// Isi parameter sesuai kebutuhan pengujian
		ID: "123", // Contoh parameter UserID
		// Isi parameter lainnya sesuai kebutuhan
	}

	expectedPaymentMethods := []adminDto.PaymentMethod{
		// Isi dengan data metode pembayaran yang diharapkan
		{
			ID: "123",
			// Isi dengan data lainnya
		},
		// Tambahkan data lain jika diperlukan
	}

	// Ekspektasi pemanggilan GetpaymentMethodByParams di repository
	mockRepo.On("GetpaymentMethodByParams", params).Return(expectedPaymentMethods, nil)

	// Pengujian
	paymentMethods, err := usecase.GetpaymentMethodByParams(params)

	// Verifikasi
	assert.NoError(t, err)
	assert.Equal(t, expectedPaymentMethods, paymentMethods)
	mockRepo.AssertExpectations(t)
}

func TestGetpaymentMethodByParams_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	params := adminDto.GetpaymentMethodParams{
		// Isi parameter sesuai kebutuhan pengujian
		ID: "123", // Contoh parameter UserID
		// Isi parameter lainnya sesuai kebutuhan
	}

	expectedError := errors.New("database error") // Contoh error yang diharapkan dari repository

	// Ekspektasi pemanggilan GetpaymentMethodByParams di repository
	mockRepo.On("GetpaymentMethodByParams", params).Return([]adminDto.PaymentMethod{}, expectedError)

	// Pengujian
	paymentMethods, err := usecase.GetpaymentMethodByParams(params)

	// Verifikasi
	assert.Error(t, err)
	assert.Nil(t, paymentMethods)
	assert.EqualError(t, err, expectedError.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetWalletByParams_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	params := adminDto.GetWalletParams{
		// Isi parameter sesuai kebutuhan pengujian
		ID: "123", // Contoh parameter UserID
		// Isi parameter lainnya sesuai kebutuhan
	}

	expectedWallets := []adminDto.Wallet{
		// Isi dengan data dompet yang diharapkan
		{
			ID: "123",
			// Isi dengan data lainnya
		},
		// Tambahkan data lain jika diperlukan
	}

	// Ekspektasi pemanggilan GetWalletByParams di repository
	mockRepo.On("GetWalletByParams", params).Return(expectedWallets, nil)

	// Pengujian
	wallets, err := usecase.GetWalletByParams(params)

	// Verifikasi
	assert.NoError(t, err)
	assert.Equal(t, expectedWallets, wallets)
	mockRepo.AssertExpectations(t)
}

func TestGetWalletByParams_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	params := adminDto.GetWalletParams{
		// Isi parameter sesuai kebutuhan pengujian
		ID: "123", // Contoh parameter UserID
		// Isi parameter lainnya sesuai kebutuhan
	}

	expectedError := errors.New("database error") // Contoh error yang diharapkan dari repository

	// Ekspektasi pemanggilan GetWalletByParams di repository
	mockRepo.On("GetWalletByParams", params).Return([]adminDto.Wallet{}, expectedError)

	// Pengujian
	wallets, err := usecase.GetWalletByParams(params)

	// Verifikasi
	assert.Error(t, err)
	assert.Nil(t, wallets)
	assert.EqualError(t, err, expectedError.Error())
	mockRepo.AssertExpectations(t)
}

func TestSavePaymentMethod_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	request := adminDto.CreatePaymentMethod{
		// Isi request sesuai kebutuhan pengujian
		PaymentName: "BCA", // Contoh nama metode pembayaran
	}

	// Ekspektasi pemanggilan SavePaymentMethod di repository
	mockRepo.On("SavePaymentMethod", mock.AnythingOfType("adminDto.PaymentMethod")).Return(nil)

	// Pengujian
	err := usecase.SavePaymentMethod(request)

	// Verifikasi
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSavePaymentMethod_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	request := adminDto.CreatePaymentMethod{
		// Isi request sesuai kebutuhan pengujian
		PaymentName: "Credit Card", // Contoh nama metode pembayaran
	}

	expectedError := errors.New("database error") // Contoh error yang diharapkan dari repository

	// Ekspektasi pemanggilan SavePaymentMethod di repository
	mockRepo.On("SavePaymentMethod", mock.AnythingOfType("adminDto.PaymentMethod")).Return(expectedError)

	// Pengujian
	err := usecase.SavePaymentMethod(request)

	// Verifikasi
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockRepo.AssertExpectations(t)
}

func TestSoftDeletePaymentMethod_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	paymentMethodID := "123" // Contoh ID metode pembayaran

	// Ekspektasi pemanggilan SoftDeletePaymentMethod di repository
	mockRepo.On("SoftDeletePaymentMethod", paymentMethodID).Return(nil)

	// Pengujian
	err := usecase.SoftDeletePaymentMethod(paymentMethodID)

	// Verifikasi
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeletePaymentMethod_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	paymentMethodID := "123" // Contoh ID metode pembayaran

	expectedError := errors.New("database error") // Contoh error yang diharapkan dari repository

	// Ekspektasi pemanggilan SoftDeletePaymentMethod di repository
	mockRepo.On("SoftDeletePaymentMethod", paymentMethodID).Return(expectedError)

	// Pengujian
	err := usecase.SoftDeletePaymentMethod(paymentMethodID)

	// Verifikasi
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdatePaymentMethod_Success(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	request := adminDto.UpdatePaymentRequest{
		ID:          "123", // Contoh ID metode pembayaran
		PaymentName: "BRI", // Contoh nama metode pembayaran yang baru
	}

	// Ekspektasi pemanggilan UpdatePaymentMethod di repository
	mockRepo.On("UpdatePaymentMethod", mock.AnythingOfType("adminDto.PaymentMethod")).Return(nil)

	// Pengujian
	err := usecase.UpdatePaymentMethod(request)

	// Verifikasi
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePaymentMethod_Failure(t *testing.T) {
	// Persiapan
	mockRepo := new(MockAdminRepository)
	usecase := adminUsecase.NewAdminUsecase(mockRepo)

	request := adminDto.UpdatePaymentRequest{
		ID:          "123",             // Contoh ID metode pembayaran
		PaymentName: "New Credit Card", // Contoh nama metode pembayaran yang baru
	}

	expectedError := errors.New("database error") // Contoh error yang diharapkan dari repository

	// Ekspektasi pemanggilan UpdatePaymentMethod di repository
	mockRepo.On("UpdatePaymentMethod", mock.AnythingOfType("adminDto.PaymentMethod")).Return(expectedError)

	// Pengujian
	err := usecase.UpdatePaymentMethod(request)

	// Verifikasi
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockRepo.AssertExpectations(t)
}
