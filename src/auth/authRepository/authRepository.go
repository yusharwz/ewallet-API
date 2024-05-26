package authRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/auth"
	"time"

	"github.com/rs/zerolog/log"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auth.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (repo *authRepository) CekEmail(email string) (resp userDto.ForgetPinResp, err error) {
	query := "SELECT email, username, pin, status FROM users WHERE email = $1 AND deleted_at IS NULL"

	err = repo.db.QueryRow(query, email).Scan(&resp.Email, &resp.Username, &resp.Unique, &resp.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Msg("email not registered")
			return userDto.ForgetPinResp{}, errors.New("email not registered")
		}
		return userDto.ForgetPinResp{}, err
	}

	return resp, nil
}

func (repo *authRepository) CekPhoneNumber(pnumber string) (resp userDto.ForgetPinResp, err error) {
	query := "SELECT email, username, pin, status FROM users WHERE phone_number = $1 AND deleted_at IS NULL"

	err = repo.db.QueryRow(query, pnumber).Scan(&resp.Email, &resp.Username, &resp.Unique, &resp.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Msg("phone number not registered")
			return userDto.ForgetPinResp{}, errors.New("phone number not registered")
		}
		return userDto.ForgetPinResp{}, err
	}

	return resp, nil
}

func (repo *authRepository) InsertCode(code, email, pnumber string) (bool, error) {
	var resp string
	var query string
	expiredCode := time.Now().Add(5 * time.Minute)

	if email != "" {
		query = "UPDATE users SET verification_code = $1, expired_code = $2 WHERE email = $3 RETURNING email;"
		if err := repo.db.QueryRow(query, code, expiredCode, email).Scan(&resp); err != nil {
			log.Error().Msg("fail Insert Code")
			return false, errors.New("fail Insert Code")
		}
	} else if pnumber != "" {
		query = "UPDATE users SET verification_code = $1, expired_code = $2 WHERE phone_number = $3 RETURNING phone_number;"
		if err := repo.db.QueryRow(query, code, expiredCode, pnumber).Scan(&resp); err != nil {
			log.Error().Msg("fail Insert Code")
			return false, errors.New("fail Insert Code")
		}
	} else {
		log.Error().Msg("both email and phone number are empty")
		return false, errors.New("both email and phone number are empty")
	}

	return true, nil
}

func (repo *authRepository) UserLogin(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {
	var expiredCode time.Time
	var count int
	queryCount := "SELECT COUNT(*) FROM users WHERE email = $1 AND deleted_at IS NULL"
	if err := repo.db.QueryRow(queryCount, req.Email).Scan(&count); err != nil {
		return resp, err
	}
	if count == 0 {
		log.Error().Msg("email not registered")
		return resp, errors.New("email not registered")
	}

	query := "SELECT id, email, pin, expired_code, roles, status FROM users WHERE email = $1 AND verification_code = $2"
	if err := repo.db.QueryRow(query, req.Email, req.Code).Scan(&resp.UserId, &resp.UserEmail, &resp.Pin, &expiredCode, &resp.Roles, &resp.Status); err != nil {
		log.Error().Msg("invalid pin or verification code")
		return resp, errors.New("invalid pin or verification code")
	}

	currentTime := time.Now().UTC().Add(7 * time.Hour)
	expiredTimestamp := expiredCode.Unix()
	currentTimestamp := currentTime.Unix()
	if currentTimestamp > expiredTimestamp {
		log.Error().Msg("verification code has expired")
		return resp, errors.New("verification code has expired")
	}

	if resp.Status != "active" {
		return resp, errors.New("account has not been activated, please check the email inbox for the activation link")
	}

	return resp, nil
}

func (repo *authRepository) UserCreate(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, unique string, err error) {

	checkUsernameQuery := "SELECT COUNT(*) FROM users WHERE username = $1"
	var usernameCount int
	if err := repo.db.QueryRow(checkUsernameQuery, req.Username).Scan(&usernameCount); err != nil {
		log.Error().Msg("failed to check username")
		return resp, "", errors.New("failed to check username")
	}
	if usernameCount > 0 {
		log.Error().Msg("username is already in use")
		return resp, "", errors.New("username is already in use")
	}

	checkEmailQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
	var emailCount int
	if err := repo.db.QueryRow(checkEmailQuery, req.Email).Scan(&emailCount); err != nil {
		log.Error().Msg("failed to check email")
		return resp, "", errors.New("failed to check email")
	}
	if emailCount > 0 {
		log.Error().Msg("email is already in use")
		return resp, "", errors.New("email is already in use")
	}

	checkPhoneQuery := "SELECT COUNT(*) FROM users WHERE phone_number = $1"
	var phoneCount int
	if err := repo.db.QueryRow(checkPhoneQuery, req.PhoneNumber).Scan(&phoneCount); err != nil {
		log.Error().Msg("failed to check phone number")
		return resp, "", errors.New("failed to check phone number")
	}
	if phoneCount > 0 {
		log.Error().Msg("phone number is already in use")
		return resp, "", errors.New("phone number is already in use")
	}

	query := `
		INSERT INTO users (fullname, username, email, pin, phone_number, roles)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, fullname, username, email, phone_number, pin
	`
	if err := repo.db.QueryRow(query, req.Fullname, req.Username, req.Email, req.Pin, req.PhoneNumber, req.Roles).Scan(&resp.Id, &resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber, &unique); err != nil {
		log.Error().Msg("fail to create user")
		return resp, "", errors.New("fail to create user")
	}

	return resp, unique, nil
}

func (repo *authRepository) UserWalletCreate(id string) (err error) {
	query := "INSERT INTO wallets (user_id) VALUES ($1)"

	if _, err := repo.db.Exec(query, id); err != nil {
		log.Error().Msg("fail to create wallet")
		return errors.New("fail to create wallet")
	}

	return nil
}

func (repo *authRepository) ActivedAccount(req userDto.ActivatedAccountReq) (err error) {

	var expiredCode time.Time
	checkCodeExpired := "SELECT expired_code  FROM users WHERE email = $1 AND username = $2 AND pin = $3 AND verification_code = $4"
	if err := repo.db.QueryRow(checkCodeExpired, req.Email, req.Fullname, req.Unique, req.Code).Scan(&expiredCode); err != nil {
		log.Error().Msg("link activation has expired")
		return errors.New("link activation has expired")
	}
	currentTime := time.Now().UTC().Add(7 * time.Hour)
	expiredTimestamp := expiredCode.Unix()
	currentTimestamp := currentTime.Unix()
	if currentTimestamp > expiredTimestamp {
		log.Error().Msg("link activation has expired")
		return errors.New("link activation has expired")
	}

	queryUpdate := `
		UPDATE users
		SET status = 'active'
		WHERE email = $1 AND username = $2 AND pin = $3 AND verification_code = $4
	`
	if _, err := repo.db.Exec(queryUpdate, req.Email, req.Fullname, req.Unique, req.Code); err != nil {
		log.Error().Msg("failed to activated your account")
		return errors.New("failed to activated your account")
	}

	return nil
}

func (repo *authRepository) SendLinkForgetPin(req userDto.ForgetPinReq) (resp userDto.ForgetPinResp, err error) {

	queryCheckEmail := `SELECT email, username, verification_code, pin FROM users WHERE email = $1 AND phone_number = $2`
	if err := repo.db.QueryRow(queryCheckEmail, req.Email, req.PhoneNumber).Scan(&resp.Email, &resp.Username, &resp.Code, &resp.Unique); err != nil {
		log.Error().Msg("invalid email or phone number")
		return userDto.ForgetPinResp{}, errors.New("invalid email or phone number")
	}
	return resp, nil
}

func (repo *authRepository) ResetPinRepo(req userDto.ForgetPinParams) error {

	var expiredCode time.Time
	checkCodeExpired := "SELECT expired_code  FROM users WHERE email = $1 AND username = $2 AND pin = $3 AND verification_code = $4"
	if err := repo.db.QueryRow(checkCodeExpired, req.Email, req.Username, req.Unique, req.Code).Scan(&expiredCode); err != nil {
		log.Error().Msg("link reset pin has expired")
		return errors.New("link reset pin has expired")
	}
	currentTime := time.Now().UTC().Add(7 * time.Hour)
	expiredTimestamp := expiredCode.Unix()
	currentTimestamp := currentTime.Unix()
	if currentTimestamp > expiredTimestamp {
		log.Error().Msg("link reset pin has expired")
		return errors.New("link reset pin has expired")
	}

	queryUpdate := `
		UPDATE users
		SET pin = $1
		WHERE email = $2 AND username = $3 AND pin = $4 AND verification_code = $5
	`

	if _, err := repo.db.Exec(queryUpdate, req.NewPin, req.Email, req.Username, req.Unique, req.Code); err != nil {
		log.Error().Msg("failed to reset pin")
		return errors.New("failed to reset pin")
	}

	return nil
}
