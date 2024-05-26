package authDelivery_test

import (
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/auth/authDelivery"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthUsecase struct{}

func (m *mockAuthUsecase) LoginCodeReqEmail(email string) error {
	if email == "error@example.com" {
		return errors.New("failed to request login code")
	}
	return nil
}

func (m *mockAuthUsecase) LoginCodeReqSMS(phoneNumber string) error {
	if phoneNumber == "08123456789" {
		return errors.New("failed to request login code via SMS")
	}
	return nil
}

func (m *mockAuthUsecase) LoginReq(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error) {
	return userDto.UserLoginResponse{
		UserId:    "user123",
		UserEmail: req.Email,
		Token:     "token123",
		Roles:     "user",
		Status:    "active",
	}, nil
}

func (m *mockAuthUsecase) CreateReq(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error) {
	return userDto.UserCreateResponse{
		Id:          "user123",
		Fullname:    req.Fullname,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}, nil
}

func (m *mockAuthUsecase) ActivatedAccount(req userDto.ActivatedAccountReq) error {
	if req.Email == "error@example.com" {
		return errors.New("failed to activate account")
	}
	return nil
}

func (m *mockAuthUsecase) ForgotPinReqUC(req userDto.ForgetPinReq) error {
	if req.Email == "error@example.com" {
		return errors.New("failed to request pin reset")
	}
	return nil
}

func (m *mockAuthUsecase) ResetPinUC(req userDto.ForgetPinParams) error {
	if req.Email == "error@example.com" {
		return errors.New("failed to reset pin")
	}
	return nil
}

func TestLoginUserCodeReuqestEmail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	group := r.Group("user")
	mockAuthUsecase := &mockAuthUsecase{}

	authDelivery.NewAuthDelivery(group, mockAuthUsecase)

	reqBody := `{"email": "user@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/request-otp/email", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginUserCodeReuqestEmail_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	group := r.Group("user")
	mockAuthUsecase := &mockAuthUsecase{}

	authDelivery.NewAuthDelivery(group, mockAuthUsecase)

	reqBody := `{"email": "error@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/request-otp/email", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
