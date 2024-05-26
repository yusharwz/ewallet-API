package adminDelivery_test

import (
	"errors"
	"final-project-enigma/model/dto/adminDto"
	adminDelivery "final-project-enigma/src/admin/adminDelivery"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAdminUsecase struct{}

func (m *mockAdminUsecase) SavePaymentMethod(req adminDto.CreatePaymentMethod) error {
	if req.PaymentName == "error" {
		return errors.New("failed to add payment method")
	}
	return nil
}

func (m *mockAdminUsecase) UpdateUser(req adminDto.UserUpdateRequest) error {
	if req.ID == "error" {
		return errors.New("failed to update user")
	}
	return nil
}

func (m *mockAdminUsecase) SoftDeletePaymentMethod(id string) error {
	if id == "error" {
		return errors.New("failed to delete payment method")
	}
	return nil
}

func (m *mockAdminUsecase) UpdatePaymentMethod(req adminDto.UpdatePaymentRequest) error {
	if req.ID == "error" {
		return errors.New("failed to update payment method")
	}
	return nil
}

func (m *mockAdminUsecase) SoftDeleteUser(id string) error {
	if id == "error" {
		return errors.New("failed to delete user")
	}
	return nil
}

func (m *mockAdminUsecase) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	if params.ID == "error" {
		return nil, errors.New("failed to get users")
	}
	return []adminDto.User{}, nil
}

func (m *mockAdminUsecase) GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	if params.ID == "error" {
		return nil, errors.New("failed to get payment methods")
	}
	return []adminDto.PaymentMethod{}, nil
}

func (m *mockAdminUsecase) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	if params.ID == "error" {
		return nil, errors.New("failed to get wallets")
	}
	return []adminDto.Wallet{}, nil
}

func (m *mockAdminUsecase) GetTransaction(page, limit int) ([]adminDto.Transaction, int, error) {
	if page == -1 {
		return nil, 0, errors.New("failed to get transactions")
	}
	return []adminDto.Transaction{}, 0, nil
}

func (m *mockAdminUsecase) GetWalletTransaction(page, limit int) ([]adminDto.WalletTransaction, int, error) {
	if page == -1 {
		return nil, 0, errors.New("failed to get wallet transactions")
	}
	return []adminDto.WalletTransaction{}, 0, nil
}

func (m *mockAdminUsecase) GetTopUpTransaction(page, limit int) ([]adminDto.TopUpTransaction, int, error) {
	if page == -1 {
		return nil, 0, errors.New("failed to get top-up transactions")
	}
	return []adminDto.TopUpTransaction{}, 0, nil
}

func TestSavePaymentMethod_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	adminUsecase := &mockAdminUsecase{}
	adminDelivery.NewAdminDelivery(r.Group("/admin"), adminUsecase)

	payload := `{"payment_name":"test"}`
	req, err := http.NewRequest(http.MethodPost, "/admin/paymentMethod", strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := `{"data":{"payment_name":"test"},"message":"succes","status_code":"01","status_message":"01"}`
	assert.Equal(t, expectedResponse, w.Body.String())
}

func TestSavePaymentMethod_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	adminUsecase := &mockAdminUsecase{}
	adminDelivery.NewAdminDelivery(r.Group("/admin"), adminUsecase)

	payload := `{"payment_name":"error"}`
	req, err := http.NewRequest(http.MethodPost, "/admin/paymentMethod", strings.NewReader(payload))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedResponse := `{"message":"failed to add payment method","status_code":"01"}`
	assert.Equal(t, expectedResponse, w.Body.String())
}
