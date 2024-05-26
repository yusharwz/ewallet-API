package authRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/userDto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCekEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewAuthRepository(db)

	mock.ExpectQuery("SELECT email, username, pin, status FROM users").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"email", "username", "pin", "status"}).
			AddRow("test@example.com", "testuser", "123456", "active"))

	resp, err := repo.CekEmail("test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "testuser", resp.Username)
	assert.Equal(t, "123456", resp.Unique)
	assert.Equal(t, "active", resp.Status)
}

func TestCekEmail_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewAuthRepository(db)

	mock.ExpectQuery("SELECT email, username, pin, status FROM users").
		WithArgs("test@example.com").
		WillReturnError(sql.ErrNoRows)

	resp, err := repo.CekEmail("test@example.com")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
	assert.Equal(t, userDto.ForgetPinResp{}, resp)
}

func TestInsertCode(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewAuthRepository(db)

	expiredCode := time.Now().Add(5 * time.Minute)
	mock.ExpectExec("UPDATE users SET verification_code = $1, expired_code = $2 WHERE email = $3 RETURNING email;").
		WithArgs("123456", expiredCode, "test@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	success, err := repo.InsertCode("123456", "test@example.com", "")

	assert.NoError(t, err)
	assert.True(t, success)
}

func TestUserLogin(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewAuthRepository(db)

	expiredCode := time.Now().Add(5 * time.Minute)
	mock.ExpectQuery("SELECT COUNT(*) FROM users").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery("SELECT id, email, pin, expired_code, roles, status FROM users").
		WithArgs("test@example.com", "123456").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "pin", "expired_code", "roles", "status"}).
			AddRow("1", "test@example.com", "123456", expiredCode, "user", "active"))

	resp, err := repo.UserLogin(userDto.UserLoginRequest{Email: "test@example.com", Pin: "123456", Code: "123456"})

	assert.NoError(t, err)
	assert.Equal(t, "1", resp.UserId)
	assert.Equal(t, "test@example.com", resp.UserEmail)
	assert.Equal(t, "123456", resp.Pin)
	assert.Equal(t, "user", resp.Roles)
	assert.Equal(t, "active", resp.Status)
}
