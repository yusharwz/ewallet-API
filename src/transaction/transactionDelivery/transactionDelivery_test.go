package transactionDelivery

import (
	"errors"
	"final-project-enigma/model/dto/transactionDtos"
	"final-project-enigma/src/transaction/transactionUseCase"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockTransactionRepository struct{}

func (m *mockTransactionRepository) GetTopUpTransaction(page int, limit int) ([]transactionDtos.TopUpTransaction, int, error) {
	if page == 1 && limit == 10 {
		return []transactionDtos.TopUpTransaction{}, 0, nil
	}
	return nil, 0, errors.New("Failed to fetch top up transactions")
}

func (m *mockTransactionRepository) GetWalletTransaction(page int, limit int) ([]transactionDtos.WalletTransaction, int, error) {
	if page == 1 && limit == 10 {
		return []transactionDtos.WalletTransaction{}, 0, nil
	}
	return nil, 0, errors.New("Failed to fetch wallet transactions")
}

func (m *mockTransactionRepository) GetTransaction(page int, limit int) ([]transactionDtos.Transaction, int, error) {
	if page == 1 && limit == 10 {
		return []transactionDtos.Transaction{}, 0, nil
	}
	return nil, 0, errors.New("Failed to fetch transactions")
}

func TestGetTopUpTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mockTransactionRepository{}
	transactionUC := transactionUseCase.NewTransactionUseCase(mockRepo)
	NewTransactionDelivery(router.Group(""), transactionUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transaction/topup?page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

func TestGetWalletTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mockTransactionRepository{}
	transactionUC := transactionUseCase.NewTransactionUseCase(mockRepo)
	NewTransactionDelivery(router.Group(""), transactionUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transaction/wallet?page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

func TestGetTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mockTransactionRepository{}
	transactionUC := transactionUseCase.NewTransactionUseCase(mockRepo)
	NewTransactionDelivery(router.Group(""), transactionUC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transaction?page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}
