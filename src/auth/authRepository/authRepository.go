package authRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/auth"
	"time"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auth.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (repo *authRepository) CekEmail(email string) (bool, error) {
	var result string
	query := "SELECT status FROM users WHERE email = $1"

	err := repo.db.QueryRow(query, email).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("email not registered")
		}
		return false, err
	}

	if result != "active" {
		return false, errors.New("account has not been activated, please check the email inbox for the activation link")
	}

	return true, nil
}

func (repo *authRepository) CekPhoneNumber(pnumber string) (bool, error) {
	var result string
	query := "SELECT status FROM users WHERE phone_number = $1"

	err := repo.db.QueryRow(query, pnumber).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("phone number not registered")
		}
		return false, err
	}

	if result != "active" {
		return false, errors.New("account has not been activated, please check the email inbox for the activation link")
	}

	return true, nil
}

func (repo *authRepository) InsertCode(code, email, pnumber string) (bool, error) {
	var resp string
	var query string
	expiredCode := time.Now().Add(5 * time.Minute)

	if email != "" {
		query = "UPDATE users SET verification_code = $1, expired_code = $2 WHERE email = $3 RETURNING email;"
		if err := repo.db.QueryRow(query, code, expiredCode, email).Scan(&resp); err != nil {
			return false, errors.New("fail Insert Code")
		}
	} else if pnumber != "" {
		query = "UPDATE users SET verification_code = $1, expired_code = $2 WHERE phone_number = $3 RETURNING phone_number;"
		if err := repo.db.QueryRow(query, code, expiredCode, pnumber).Scan(&resp); err != nil {
			return false, errors.New("fail Insert Code")
		}
	} else {
		return false, errors.New("both email and phone number are empty")
	}

	return true, nil
}

func (repo *authRepository) UserLogin(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {
	var expiredCode time.Time

	query := "SELECT id, email, pin, expired_code, roles, status FROM users WHERE phone_number = $1 AND verification_code = $2"
	if err := repo.db.QueryRow(query, req.PhoneNumber, req.Code).Scan(&resp.UserId, &resp.UserEmail, &resp.Pin, &expiredCode, &resp.Roles, &resp.Status); err != nil {
		return resp, errors.New("invalid pin or verification code")
	}

	currentTime := time.Now().UTC().Add(7 * time.Hour)
	expiredTimestamp := expiredCode.Unix()
	currentTimestamp := currentTime.Unix()
	if currentTimestamp > expiredTimestamp {
		return resp, errors.New("verification code has expired")
	}

	if resp.Status != "active" {
		return resp, errors.New("account has not been activated, please check the email inbox for the activation link")
	}

	return resp, nil
}

func (repo *authRepository) UserCreate(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, unique string, err error) {

	checkEmailQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
	var emailCount int
	if err := repo.db.QueryRow(checkEmailQuery, req.Email).Scan(&emailCount); err != nil {
		return resp, "", errors.New("failed to check email")
	}
	if emailCount > 0 {
		return resp, "", errors.New("email is already in use")
	}

	checkUsernameQuery := "SELECT COUNT(*) FROM users WHERE username = $1"
	var usernameCount int
	if err := repo.db.QueryRow(checkUsernameQuery, req.Username).Scan(&usernameCount); err != nil {
		return resp, "", errors.New("failed to check username")
	}
	if usernameCount > 0 {
		return resp, "", errors.New("username is already in use")
	}

	checkPhoneQuery := "SELECT COUNT(*) FROM users WHERE phone_number = $1"
	var phoneCount int
	if err := repo.db.QueryRow(checkPhoneQuery, req.PhoneNumber).Scan(&phoneCount); err != nil {
		return resp, "", errors.New("failed to check phone number")
	}
	if phoneCount > 0 {
		return resp, "", errors.New("phone number is already in use")
	}

	query := `
		INSERT INTO users (fullname, username, email, pin, phone_number, roles)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, fullname, username, email, phone_number, pin
	`
	if err := repo.db.QueryRow(query, req.Fullname, req.Username, req.Email, req.Pin, req.PhoneNumber, req.Roles).Scan(&resp.Id, &resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber, &unique); err != nil {
		return resp, "", errors.New("fail to create user")
	}

	return resp, unique, nil
}

func (repo *authRepository) UserWalletCreate(id string) (err error) {
	query := "INSERT INTO wallets (user_id) VALUES ($1)"

	if _, err := repo.db.Exec(query, id); err != nil {
		return errors.New("fail to create wallet")
	}

	return nil
}

func (repo *authRepository) ActivedAccount(req userDto.ActivatedAccountReq) (err error) {

	var expiredCode time.Time
	checkCodeExpired := "SELECT expired_code  FROM users WHERE email = $1 AND username = $2 AND pin = $3 AND verification_code = $4"
	if err := repo.db.QueryRow(checkCodeExpired, req.Email, req.Fullname, req.Unique, req.Code).Scan(&expiredCode); err != nil {
		return errors.New("link activation has expired")
	}
	currentTime := time.Now().UTC().Add(7 * time.Hour)
	expiredTimestamp := expiredCode.Unix()
	currentTimestamp := currentTime.Unix()
	if currentTimestamp > expiredTimestamp {
		return errors.New("link activation has expired")
	}

	queryUpdate := `
		UPDATE users
		SET status = 'active'
		WHERE email = $1 AND username = $2 AND pin = $3 AND verification_code = $4
	`
	if _, err := repo.db.Exec(queryUpdate, req.Email, req.Fullname, req.Unique, req.Code); err != nil {
		return errors.New("failed to activated your account")
	}

	return nil
}
