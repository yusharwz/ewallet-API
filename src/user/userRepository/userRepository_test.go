package userRepository_test

import (
	"bytes"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/user/userRepository"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestEditUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := resty.New()
	repo := userRepository.NewUserRepository(db, client)

	req := userDto.UserUpdateReq{
		UserId:      "1",
		Fullname:    "John Doe",
		Username:    "johndoe",
		Email:       "john@example.com",
		PhoneNumber: "1234567890",
	}

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE username = \\$1 AND id != \\$2").
		WithArgs(req.Username, req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 AND id != \\$2").
		WithArgs(req.Email, req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE phone_number = \\$1 AND id != \\$2").
		WithArgs(req.PhoneNumber, req.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec("UPDATE users SET fullname = \\$1, username = \\$2, email = \\$3, phone_number = \\$4 WHERE id = \\$5").
		WithArgs(req.Fullname, req.Username, req.Email, req.PhoneNumber, req.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.EditUserData(req)
	assert.NoError(t, err)
}

func TestUserUploadImage(t *testing.T) {
	// Membuat file sementara untuk pengujian
	tempFile, err := os.CreateTemp("", "test*.jpg")
	if err != nil {
		t.Fatalf("An error '%s' was not expected when creating a temp file", err)
	}
	defer os.Remove(tempFile.Name()) // Menghapus file sementara setelah pengujian selesai

	// Menulis data ke file sementara
	_, err = tempFile.Write([]byte("test image content"))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when writing to the temp file", err)
	}
	tempFile.Close()

	// Membuat multipart form file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile.Name()))
	if err != nil {
		t.Fatalf("An error '%s' was not expected when creating multipart form file", err)
	}
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening the temp file", err)
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("An error '%s' was not expected when copying file content to form part", err)
	}
	writer.Close()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mocking SQL Expectations
	// Contoh expectation jika UserUploadImage berinteraksi dengan database
	mock.ExpectExec("INSERT INTO images").WillReturnResult(sqlmock.NewResult(1, 1))

	client := resty.New()
	repo := userRepository.NewUserRepository(db, client)
	if repo == nil {
		t.Fatalf("Failed to create a new userRepository instance")
	}

	req := userDto.UploadImagesRequest{
		File: file,
	}

	resp, err := repo.UserUploadImage(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Url)
}

func TestImageToDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := resty.New()
	repo := userRepository.NewUserRepository(db, client)

	userId := "1"
	req := userDto.UploadImagesResponse{
		Url: "http://example.com/image.jpg",
	}

	mock.ExpectExec("UPDATE users SET image_url = \\$1, updated_at = \\$2 WHERE id = \\$3 AND deleted_at IS NULL").
		WithArgs(req.Url, sqlmock.AnyArg(), userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.ImageToDB(userId, req)
	assert.NoError(t, err)
}

func TestGetDataUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := resty.New()
	repo := userRepository.NewUserRepository(db, client)

	userId := "1"
	expectedResponse := userDto.UserGetDataResponse{
		Fullname:     "John Doe",
		Username:     "johndoe",
		Email:        "john@example.com",
		PhoneNumber:  "1234567890",
		ProfilImages: "http://example.com/image.jpg",
	}

	mock.ExpectQuery("SELECT fullname, username, email, phone_number, image_url FROM users WHERE id = \\$1 AND deleted_at IS NULL").
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"fullname", "username", "email", "phone_number", "image_url"}).
			AddRow(expectedResponse.Fullname, expectedResponse.Username, expectedResponse.Email, expectedResponse.PhoneNumber, expectedResponse.ProfilImages))

	resp, err := repo.GetDataUserRepo(userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestGetBalanceInfoRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := resty.New()
	repo := userRepository.NewUserRepository(db, client)

	userId := "1"
	expectedBalance := 100.0

	mock.ExpectQuery("SELECT balance FROM wallets WHERE user_id = \\$1 AND deleted_at IS NULL").
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(fmt.Sprintf("%.2f", expectedBalance)))

	resp, err := repo.GetBalanceInfoRepo(userId)
	assert.NoError(t, err)

	// Convert the balance from string to float64
	balance, err := strconv.ParseFloat(resp.Balance, 64)
	if err != nil {
		t.Fatalf("An error '%s' was not expected when converting balance to float64", err)
	}

	assert.Equal(t, expectedBalance, balance)
}

func TestGetTransactionRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := resty.New()
	repo := userRepository.NewUserRepository(db, client)

	params := userDto.GetTransactionParams{
		UserId:       "1",
		TrxId:        "",
		TrxDateStart: "",
		TrxDateEnd:   "",
		TrxStatus:    "",
		Page:         "1",
		Limit:        "10",
	}

	expectedTransaction := userDto.GetTransactionResponse{
		TransactionId:   "1",
		TransactionType: "",
		Amount:          "100.0",
		Description:     "Test transaction",
		TransactionDate: time.Now().Format(time.RFC3339),
		Status:          "success",
	}

	mock.ExpectQuery("SELECT id, amount, description, created_at, status FROM \\(SELECT t.id, t.amount, t.description, t.created_at, t.status FROM transactions t WHERE t.user_id = \\$1 UNION SELECT t.id, t.amount, t.description, t.created_at, t.status FROM transactions t JOIN wallet_transactions wt ON t.id = wt.transaction_id JOIN wallets w ON wt.from_wallet_id = w.id OR wt.to_wallet_id = w.id WHERE w.user_id = \\$1\\) sub WHERE 1=1 LIMIT 10 OFFSET 0").
		WithArgs(params.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "description", "created_at", "status"}).
			AddRow(expectedTransaction.TransactionId, expectedTransaction.Amount, expectedTransaction.Description, expectedTransaction.TransactionDate, expectedTransaction.Status))

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM \\(SELECT t.id, t.amount, t.description, t.created_at, t.status FROM transactions t WHERE t.user_id = \\$1 UNION SELECT t.id, t.amount, t.description, t.created_at, t.status FROM transactions t JOIN wallet_transactions wt ON t.id = wt.transaction_id JOIN wallets w ON wt.from_wallet_id = w.id OR wt.to_wallet_id = w.id WHERE w.user_id = \\$1\\) sub WHERE 1=1").
		WithArgs(params.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	resp, total, err := repo.GetTransactionRepo(params)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, expectedTransaction.TransactionId, resp[0].TransactionId)
	assert.Equal(t, expectedTransaction.Amount, resp[0].Amount)
	assert.Equal(t, expectedTransaction.Description, resp[0].Description)
	assert.Equal(t, expectedTransaction.Status, resp[0].Status)
	assert.Equal(t, 1, total)
}
