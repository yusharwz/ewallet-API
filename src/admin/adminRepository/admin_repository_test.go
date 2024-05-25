package adminRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/adminDto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersByParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("success - get users by params", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "fullname", "username", "image_url", "pin", "email", "phone_number", "roles", "status", "created_at"}).
			AddRow("1", "John Doe", "johndoe", "url", "1234", "johndoe@example.com", "123456789", "admin", "active", time.Now())
		mock.ExpectQuery("SELECT id, fullname, username, image_url, pin, email, phone_number, roles, status, created_at FROM users WHERE deleted_at IS NULL").
			WillReturnRows(rows)

		params := adminDto.GetUserParams{}
		users, err := repo.GetUsersByParams(params)

		assert.NoError(t, err)
		assert.Len(t, users, 1)
	})

	t.Run("failure - invalid start date format", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, fullname, username, image_url, pin, email, phone_number, roles, status, created_at FROM users WHERE deleted_at IS NULL").
			WillReturnError(errors.New("query error"))

		params := adminDto.GetUserParams{StartDate: "invalid-date"}
		_, err := repo.GetUsersByParams(params)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
	})

	t.Run("failure - query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, fullname, username, image_url, pin, email, phone_number, roles, status, created_at FROM users WHERE deleted_at IS NULL").
			WillReturnError(errors.New("query error"))

		params := adminDto.GetUserParams{}
		_, err := repo.GetUsersByParams(params)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
	})
}

func TestSoftDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("success - soft delete user", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET deleted_at").
			WithArgs(sqlmock.AnyArg(), "1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SoftDeleteUser("1")
		assert.NoError(t, err)
	})

	t.Run("failure - no rows affected", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET deleted_at").
			WithArgs(sqlmock.AnyArg(), "1").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.SoftDeleteUser("1")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("failure - exec error", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET deleted_at").
			WithArgs(sqlmock.AnyArg(), "1").
			WillReturnError(errors.New("exec error"))

		err := repo.SoftDeleteUser("1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exec error")
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("success - update user", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id =").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE username =").
			WithArgs("johndoe", "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email =").
			WithArgs("johndoe@example.com", "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE phone_number =").
			WithArgs("123456789", "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec("UPDATE users SET fullname =").
			WithArgs("John Doe", "johndoe", "johndoe@example.com", "123456789", "1234", sqlmock.AnyArg(), "1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		user := adminDto.User{ID: "1", Fullname: "John Doe", Username: "johndoe", Email: "johndoe@example.com", PhoneNumber: "123456789", Pin: "1234"}
		err := repo.UpdateUser(user)
		assert.NoError(t, err)
	})

	t.Run("failure - user does not exist", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id =").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		user := adminDto.User{ID: "1"}
		err := repo.UpdateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "user does not exist", err.Error())
	})

	t.Run("failure - username exists for another user", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id =").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE username =").
			WithArgs("johndoe", "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		user := adminDto.User{ID: "1", Username: "johndoe"}
		err := repo.UpdateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "username already exists for another user", err.Error())
	})
}

func TestSavePaymentMethod(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("success - save payment method", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER").
			WithArgs("Credit Card").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectExec("INSERT INTO payment_method\\(payment_name\\) VALUES").
			WithArgs("Credit Card").
			WillReturnResult(sqlmock.NewResult(1, 1))

		paymentMethod := adminDto.PaymentMethod{PaymentName: "Credit Card"}
		err := repo.SavePaymentMethod(paymentMethod)
		assert.NoError(t, err)
	})

	t.Run("failure - payment method name exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM payment_method WHERE LOWER\\(payment_name\\) = LOWER").
			WithArgs("Credit Card").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		paymentMethod := adminDto.PaymentMethod{PaymentName: "Credit Card"}
		err := repo.SavePaymentMethod(paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, "payment method name already exists", err.Error())
	})
}
func TestGetPaymentMethodByParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("success - get payment methods by params", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "payment_name", "created_at"}).
			AddRow("1", "Credit Card", time.Now())
		mock.ExpectQuery("^SELECT id, payment_name,created_at FROM payment_method WHERE 1=1 AND deleted_at IS NULL$").
			WillReturnRows(rows)

		params := adminDto.GetPaymentMethodParams{}
		paymentMethods, err := repo.GetpaymentMethodByParams(params)

		assert.NoError(t, err)
		assert.Len(t, paymentMethods, 1)
	})

}

func TestGetWalletByParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("success - get wallets by params", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "fullname", "username"}).
			AddRow("1", "1", 100.00, time.Now(), "John Doe", "johndoe")
		mock.ExpectQuery("SELECT w.id, w.user_id, w.balance, w.created_at, u.fullname, u.username FROM wallets w JOIN users u").
			WillReturnRows(rows)

		params := adminDto.GetWalletParams{}
		wallets, err := repo.GetWalletByParams(params)

		assert.NoError(t, err)
		assert.Len(t, wallets, 1)
	})

}
