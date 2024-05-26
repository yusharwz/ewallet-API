package adminDelivery

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"final-project-enigma/model/dto/adminDto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAdminUsecase struct{}

func (m *mockAdminUsecase) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	return []adminDto.User{}, nil
}

func (m *mockAdminUsecase) SoftDeleteUser(userID string) error {
	return nil
}

func (m *mockAdminUsecase) UpdateUser(user adminDto.UserUpdateRequest) error {
	return nil
}

func (m *mockAdminUsecase) GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	return []adminDto.PaymentMethod{}, nil
}

func (m *mockAdminUsecase) SavePaymentMethod(paymentMethod adminDto.CreatePaymentMethod) error {
	if paymentMethod.PaymentName == "existing" {
		return errors.New("payment method already exists")
	}
	return nil
}

func (m *mockAdminUsecase) SoftDeletePaymentMethod(paymentMethodID string) error {
	return nil
}

func (m *mockAdminUsecase) UpdatePaymentMethod(paymentMethod adminDto.UpdatePaymentRequest) error {
	return nil
}

func (m *mockAdminUsecase) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	return []adminDto.Wallet{}, nil
}

func TestSavePaymentMethod_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.POST("/admin/paymentMethod", handler.SavePaymentMethod)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"payment_name": "oyo"}`
		req, _ := http.NewRequest("POST", "/admin/paymentMethod", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failure - payment method already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"payment_name": "existing"}`
		req, _ := http.NewRequest("POST", "/admin/paymentMethod", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "payment method already exists")
	})
}

func TestSavePaymentMethod_failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.POST("/admin/paymentMethod", handler.SavePaymentMethod)

	t.Run("failure - empty payment name", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"payment_name": ""}`
		req, _ := http.NewRequest("POST", "/admin/paymentMethod", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "required")
	})
}
func TestUpdateUser_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.PUT("/admin/user/:id", handler.UpdateUser)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"id": "1", "fullname": "John Doe", "username": "johndoe", "email": "johndoe@example.com", "phone_number": "0813456789", "pin": "080402"}`
		req, _ := http.NewRequest("PUT", "/admin/user/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSoftDeleteUser_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.DELETE("/admin/user/:id", handler.SoftDeleteUser)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/admin/user/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestSoftDeletePaymentMethod_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.DELETE("/admin/paymentMethod/:id", handler.SoftDeletePaymentMethod)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/admin/paymentMethod/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdatePaymentMethod_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.PUT("/admin/paymentMethod/:id", handler.UpdatePaymentMethod)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"payment_name": "New Payment Method"}`
		req, _ := http.NewRequest("PUT", "/admin/paymentMethod/1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetUsersByParams_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.GET("/admin/users", handler.GetUsersByParams)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/users?id=1&fullname=John&username=johndoe&email=johndoe@example.com&phoneNumber=123456789&roles=admin&status=active&startDate=2024-01-01&endDate=2024-05-01&page=1&limit=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetpaymentMethodByParams_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.GET("/admin/paymentMethod", handler.GetpaymentMethodByParams)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/paymentMethod?id=1&payment_name=credit&created_at=2024-01-01&page=1&limit=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetWalletByParams_success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	adminUsecase := &mockAdminUsecase{}
	handler := &adminDelivery{adminUsecase: adminUsecase}
	router.GET("/admin/wallet", handler.GetWalletByParams)

	t.Run("success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/wallet?id=1&user_id=1&fullname=John&username=johndoe&created_at=2024-01-01&min_balance=100&max_balance=1000&page=1&limit=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
