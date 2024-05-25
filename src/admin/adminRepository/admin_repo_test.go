package adminRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/adminDto"
	"testing"
	"time"

	// "time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersByParams(t *testing.T) {
	// Inisialisasi database mock dan repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully get users by params", func(t *testing.T) {
		params := adminDto.GetUserParams{
			ID: "123",
		}

		mock.ExpectQuery("SELECT id, fullname, username, image_url, pin, email, phone_number, roles, status, created_at FROM users WHERE deleted_at IS NULL AND id = \\$1").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "fullname", "username", "image_url", "pin", "email", "phone_number", "roles", "status", "created_at"}).
				AddRow("123", "John Doe", "johndoe", "http://example.com/avatar.jpg", "1234", "johndoe@example.com", "123456789", "user", "active", time.Now()))

		users, err := repo.GetUsersByParams(params)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
	})

	t.Run("User not found", func(t *testing.T) {
		params := adminDto.GetUserParams{
			Username: "johndoe",
		}

		mock.ExpectQuery("SELECT id, fullname, username, image_url, pin, email, phone_number, roles, status, created_at FROM users WHERE deleted_at IS NULL AND username = \\$1").
			WithArgs("johndoe").
			WillReturnRows(sqlmock.NewRows([]string{}))

		_, err := repo.GetUsersByParams(params)
		assert.Error(t, err)
		assert.Equal(t, errors.New("user with username 'johndoe' not found"), err)
	})

	// Lakukan pengujian tambahan untuk skenario lainnya (mis. email tidak ditemukan, nomor telepon tidak ditemukan, dsb.)
}

func TestSoftDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully soft delete a user", func(t *testing.T) {
		userID := "123"
		query := "UPDATE users SET deleted_at=\\$1 WHERE id=\\$2 AND deleted_at IS NULL"

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SoftDeleteUser(userID)
		assert.NoError(t, err)
	})

	t.Run("User not found or already deleted", func(t *testing.T) {
		userID := "123"
		query := "UPDATE users SET deleted_at=\\$1 WHERE id=\\$2 AND deleted_at IS NULL"

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.SoftDeleteUser(userID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Database error", func(t *testing.T) {
		userID := "123"
		query := "UPDATE users SET deleted_at=\\$1 WHERE id=\\$2 AND deleted_at IS NULL"

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), userID).
			WillReturnError(errors.New("db error"))

		err := repo.SoftDeleteUser(userID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestUpdateUser(t *testing.T) {
	// Inisialisasi database mock dan repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully update user", func(t *testing.T) {
		user := adminDto.User{
			ID:           "123",
			Fullname:     "John Doe",
			Username:     "johndoe",
			Email:        "johndoe@example.com",
			PhoneNumber:  "123456789",
			Pin:          "1234",
			UpdatedAt:    time.Now(),
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id = \\$1 AND deleted_at IS NULL\\)").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE username = \\$1 AND id <> \\$2 AND deleted_at IS NULL\\)").
			WithArgs("johndoe", "123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\$1 AND id <> \\$2 AND deleted_at IS NULL\\)").
			WithArgs("johndoe@example.com", "123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE phone_number = \\$1 AND id <> \\$2 AND deleted_at IS NULL\\)").
			WithArgs("123456789", "123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("UPDATE users SET fullname = \\$1, username = \\$2, email = \\$3, phone_number = \\$4, pin = \\$5, updated_at = \\$6 WHERE id = \\$7 AND deleted_at IS NULL").
			WithArgs("John Doe", "johndoe", "johndoe@example.com", "123456789", "1234", sqlmock.AnyArg(), "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUser(user)
		assert.NoError(t, err)
	})

	t.Run("User does not exist", func(t *testing.T) {
		user := adminDto.User{
			ID: "123",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id = \\$1 AND deleted_at IS NULL\\)").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		err := repo.UpdateUser(user)
		assert.Error(t, err)
		assert.Equal(t, errors.New("user does not exist"), err)
	})

	t.Run("Username already exists for another user", func(t *testing.T) {
		user := adminDto.User{
			ID:       "123",
			Username: "johndoe",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id = \\$1 AND deleted_at IS NULL\\)").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE username = \\$1 AND id <> \\$2 AND deleted_at IS NULL\\)").
			WithArgs("johndoe", "123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.UpdateUser(user)
		assert.Error(t, err)
		assert.Equal(t, errors.New("username already exists for another user"), err)
	})

	// Lakukan pengujian tambahan untuk skenario lainnya (mis. email sudah ada, nomor telepon sudah ada, dsb.)
}

func TestGetpaymentMethodByParams(t *testing.T) {
	// Inisialisasi database mock dan repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully get payment methods by params", func(t *testing.T) {
		params := adminDto.GetpaymentMethodParams{
			ID: "123",
		}

		mock.ExpectQuery("SELECT id, payment_name,created_at FROM payment_method WHERE 1=1 AND deleted_at IS NULL AND id = \\$1").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "payment_name", "created_at"}).
				AddRow("123", "Credit Card", time.Now()))

		paymentMethods, err := repo.GetpaymentMethodByParams(params)
		assert.NoError(t, err)
		assert.Len(t, paymentMethods, 1)
	})

	t.Run("Payment method not found", func(t *testing.T) {
		params := adminDto.GetpaymentMethodParams{
			PaymentName: "Credit Card",
		}

		mock.ExpectQuery("SELECT id, payment_name,created_at FROM payment_method WHERE 1=1 AND deleted_at IS NULL AND payment_name LIKE \\$1").
			WithArgs("%Credit Card%").
			WillReturnRows(sqlmock.NewRows([]string{}))

		_, err := repo.GetpaymentMethodByParams(params)
		assert.Error(t, err)
		assert.Equal(t, errors.New("payment with name 'Credit Card' not found"), err)
	})

	// Lakukan pengujian tambahan untuk skenario lainnya (mis. ID tidak ditemukan, dsb.)
}

func TestSavePaymentMethod(t *testing.T) {
	// Inisialisasi database mock dan repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully save new payment method", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			PaymentName: "Credit Card",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Credit Card").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO payment_method\\(payment_name\\) VALUES\\(\\$1\\)").
			WithArgs("Credit Card").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SavePaymentMethod(paymentMethod)
		assert.NoError(t, err)
	})

	t.Run("Payment method already exists", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			PaymentName: "PayPal",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("PayPal").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.SavePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, errors.New("payment method name already exists"), err)
	})

	t.Run("Error checking if payment method exists", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			PaymentName: "Debit Card",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Debit Card").
			WillReturnError(errors.New("database error"))

		err := repo.SavePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, errors.New("database error"), err)
	})

	t.Run("Error saving payment method", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			PaymentName: "Bank Transfer",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Bank Transfer").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO payment_method\\(payment_name\\) VALUES\\(\\$1\\)").
			WithArgs("Bank Transfer").
			WillReturnError(errors.New("database error"))

		err := repo.SavePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, errors.New("database error"), err)
	})
}

func TestSoftDeletePaymentMethod(t *testing.T) {
	// Inisialisasi database mock dan repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully soft delete a payment method", func(t *testing.T) {
		paymentMethodID := "123"
		query := "UPDATE payment_method SET deleted_at=\\$1 WHERE id=\\$2 AND deleted_at IS NULL"

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), paymentMethodID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SoftDeletePaymentMethod(paymentMethodID)
		assert.NoError(t, err)
	})

	t.Run("Payment method not found or already deleted", func(t *testing.T) {
		paymentMethodID := "999"
		query := "UPDATE payment_method SET deleted_at=\\$1 WHERE id=\\$2 AND deleted_at IS NULL"

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), paymentMethodID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.SoftDeletePaymentMethod(paymentMethodID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Database error during soft delete", func(t *testing.T) {
		paymentMethodID := "123"
		query := "UPDATE payment_method SET deleted_at=\\$1 WHERE id=\\$2 AND deleted_at IS NULL"

		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), paymentMethodID).
			WillReturnError(errors.New("database error"))

		err := repo.SoftDeletePaymentMethod(paymentMethodID)
		assert.Error(t, err)
		assert.Equal(t, errors.New("database error"), err)
	})
}

func TestUpdatePaymentMethod(t *testing.T) {
	// Inisialisasi database mock dan repository
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing sqlmock: %v", err)
	}
	defer db.Close()

	repo := &adminRepo{db}

	t.Run("Successfully update payment method", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			ID:          "123",
			PaymentName: "Credit Card",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Credit Card").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("UPDATE payment_method SET payment_name=\\$1, updated_at=\\$2 WHERE id=\\$3 AND deleted_at IS NULL").
			WithArgs("Credit Card", sqlmock.AnyArg(), "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdatePaymentMethod(paymentMethod)
		assert.NoError(t, err)
	})

	t.Run("Payment method name already exists", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			ID:          "123",
			PaymentName: "PayPal",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("PayPal").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.UpdatePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, errors.New("payment method name already exists"), err)
	})

	t.Run("Error checking payment method existence", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			ID:          "123",
			PaymentName: "Debit Card",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Debit Card").
			WillReturnError(errors.New("kesalahan db"))

		err := repo.UpdatePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, errors.New("kesalahan db"), err)
	})

	t.Run("Error updating payment method", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			ID:          "123",
			PaymentName: "Bank Transfer",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Bank Transfer").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("UPDATE payment_method SET payment_name=\\$1, updated_at=\\$2 WHERE id=\\$3 AND deleted_at IS NULL").
			WithArgs("Bank Transfer", sqlmock.AnyArg(), "123").
			WillReturnError(errors.New("kesalahan db"))

		err := repo.UpdatePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, errors.New("kesalahan db"), err)
	})

	t.Run("Payment method not found or already deleted", func(t *testing.T) {
		paymentMethod := adminDto.PaymentMethod{
			ID:          "999",
			PaymentName: "Bitcoin",
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER\\(\\$1\\) AND deleted_at IS NULL\\)").
			WithArgs("Bitcoin").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("UPDATE payment_method SET payment_name=\\$1, updated_at=\\$2 WHERE id=\\$3 AND deleted_at IS NULL").
			WithArgs("Bitcoin", sqlmock.AnyArg(), "999").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdatePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

