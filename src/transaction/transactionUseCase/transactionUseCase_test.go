package transactionUseCase_test

import (
	"errors"
	"final-project-enigma/model/dto/transactionDtos"
	"final-project-enigma/src/transaction"
	"final-project-enigma/src/transaction/transactionUseCase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTransactionUseCase struct {
	mock.Mock
}

func (m *MockTransactionUseCase) GetTopUpTransaction(page int, limit int) ([]transactionDtos.TopUpTransaction, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]transactionDtos.TopUpTransaction), args.Int(1), args.Error(2)
}

func (m *MockTransactionUseCase) GetWalletTransaction(page int, limit int) ([]transactionDtos.WalletTransaction, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]transactionDtos.WalletTransaction), args.Int(1), args.Error(2)
}

func (m *MockTransactionUseCase) GetTransaction(page int, limit int) ([]transactionDtos.Transaction, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]transactionDtos.Transaction), args.Int(1), args.Error(2)
}

type TransactionUseCaseTestSuite struct {
	suite.Suite
	mockRepo *MockTransactionUseCase
	useCase  transaction.TransactionUseCase
}

func (suite *TransactionUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockTransactionUseCase)
	suite.useCase = transactionUseCase.NewTransactionUseCase(suite.mockRepo)
}

func (suite *TransactionUseCaseTestSuite) TestGetWalletTransaction_Success() {
	// Prepare test data
	page := 1
	limit := 10
	expectedWalletTransactions := []transactionDtos.WalletTransaction{
		{Id: "1", TransactionId: "1", FromWalletId: "1", ToWalletId: "2", Created_at: time.Now()},
		{Id: "2", TransactionId: "2", FromWalletId: "1", ToWalletId: "2", Created_at: time.Now()},
	}
	expectedTotalCount := len(expectedWalletTransactions)

	// Mock the repository method
	suite.mockRepo.On("GetWalletTransaction", page, limit).Return(expectedWalletTransactions, expectedTotalCount, nil)

	// Call the method under test
	walletTransactions, totalCount, err := suite.useCase.GetWalletTransaction(page, limit)

	// Check the results
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), walletTransactions)
	assert.Equal(suite.T(), expectedTotalCount, totalCount)
	assert.Equal(suite.T(), expectedWalletTransactions, walletTransactions)

	// Ensure all expectations on the repository were met
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransaction_Success() {
	// Prepare test data
	page := 1
	limit := 10
	expectedTransactions := []transactionDtos.Transaction{
		{Id: "1", TransactionType: "type1", Amount: 100, Created_at: time.Now()},
		{Id: "2", TransactionType: "type2", Amount: 200, Created_at: time.Now()},
	}
	expectedTotalCount := len(expectedTransactions)

	// Mock the repository method
	suite.mockRepo.On("GetTransaction", page, limit).Return(expectedTransactions, expectedTotalCount, nil)

	// Call the method under test
	transactions, totalCount, err := suite.useCase.GetTransaction(page, limit)

	// Check the results
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transactions)
	assert.Equal(suite.T(), expectedTotalCount, totalCount)
	assert.Equal(suite.T(), expectedTransactions, transactions)

	// Ensure all expectations on the repository were met
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTopUpTransaction_Success() {
	// Prepare test data
	page := 1
	limit := 10
	expectedTopUpTransactions := []transactionDtos.TopUpTransaction{
		{Id: "1", TransactionId: "1", PaymentMethodId: "1", Created_at: time.Now()},
		{Id: "2", TransactionId: "2", PaymentMethodId: "1", Created_at: time.Now()},
	}
	expectedTotalCount := len(expectedTopUpTransactions)

	suite.mockRepo.On("GetTopUpTransaction", page, limit).Return(expectedTopUpTransactions, expectedTotalCount, nil)

	// Call the method under test
	topUpTransactions, totalCount, err := suite.useCase.GetTopUpTransaction(page, limit)

	// Check the results
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), topUpTransactions)
	assert.Equal(suite.T(), expectedTotalCount, totalCount)
	assert.Equal(suite.T(), expectedTopUpTransactions, topUpTransactions)

	// Ensure all expectations on the repository were met
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetWalletTransaction_Error() {
	// Prepare test data
	page := 1
	limit := 10
	expectedError := errors.New("repository error")

	// Mock the repository method
	suite.mockRepo.On("GetWalletTransaction", page, limit).Return([]transactionDtos.WalletTransaction{}, 0, expectedError)

	// Call the method under test
	walletTransactions, totalCount, err := suite.useCase.GetWalletTransaction(page, limit)

	// Check the results
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), walletTransactions)
	assert.Equal(suite.T(), 0, totalCount)
	assert.Equal(suite.T(), expectedError, err)

	// Ensure all expectations on the repository were met
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTransaction_Error() {
	// Prepare test data
	page := 1
	limit := 10
	expectedError := errors.New("repository error")

	// Mock the repository method to return the correct types and an error
	suite.mockRepo.On("GetTransaction", page, limit).Return([]transactionDtos.Transaction{}, 0, expectedError)

	// Call the method under test
	transactions, totalCount, err := suite.useCase.GetTransaction(page, limit)

	// Check the results
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), transactions)
	assert.Equal(suite.T(), 0, totalCount)
	assert.Equal(suite.T(), expectedError, err)

	// Ensure all expectations on the repository were met
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetTopUpTransaction_Error() {
	// Prepare test data
	page := 1
	limit := 10
	expectedError := errors.New("repository error")

	// Mock the repository method
	suite.mockRepo.On("GetTopUpTransaction", page, limit).Return([]transactionDtos.TopUpTransaction{}, 0, expectedError)

	// Call the method under test
	topUpTransactions, totalCount, err := suite.useCase.GetTopUpTransaction(page, limit)

	// Check the results
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), topUpTransactions)
	assert.Equal(suite.T(), 0, totalCount)
	assert.Equal(suite.T(), expectedError, err)

	// Ensure all expectations on the repository were met
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestTransactionUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseTestSuite))
}
