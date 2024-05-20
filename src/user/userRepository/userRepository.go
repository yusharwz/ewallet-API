package userRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/user"
	"fmt"
	"strconv"
	"time"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) CekEmail(email string) (bool, error) {
	var result int
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	if err := repo.db.QueryRow(query, email).Scan(&result); err != nil {
		return false, errors.New("email not found")
	}

	if result == 0 {
		return false, errors.New("email not found")
	}

	return true, nil
}

func (repo *userRepository) CekPhoneNumber(pnumber string) (bool, error) {
	var result int
	query := "SELECT COUNT(*) FROM users WHERE phone_number = $1"
	if err := repo.db.QueryRow(query, pnumber).Scan(&result); err != nil {
		return false, errors.New("phone number not found")
	}

	if result == 0 {
		return false, errors.New("phone number not found")
	}

	return true, nil
}

func (repo *userRepository) InsertCode(code, email, pnumber string) (bool, error) {
	var resp string
	var query string
	expiredCode := time.Now().Add(5 * time.Minute)

	if email != "" {
		query = "UPDATE users SET verification_code = $1, expired_code = $2 WHERE email = $3 RETURNING email;"
		if err := repo.db.QueryRow(query, code, expiredCode, email).Scan(&resp); err != nil {
			fmt.Println(err)
			return false, errors.New("fail Insert Code")
		}
	} else if pnumber != "" {
		query = "UPDATE users SET verification_code = $1, expired_code = $2 WHERE phone_number = $3 RETURNING phone_number;"
		if err := repo.db.QueryRow(query, code, expiredCode, pnumber).Scan(&resp); err != nil {
			fmt.Println(err)
			return false, errors.New("fail Insert Code")
		}
	} else {
		return false, errors.New("both email and phone number are empty")
	}

	return true, nil
}

func (repo *userRepository) UserLogin(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {
	var expiredCode time.Time

	query := "SELECT id, email, pin, expired_code FROM users WHERE phone_number = $1 AND verification_code = $2"
	if err := repo.db.QueryRow(query, req.PhoneNumber, req.Code).Scan(&resp.UserId, &resp.UserEmail, &resp.Pin, &expiredCode); err != nil {
		fmt.Println(req.Code)
		return resp, errors.New("invalid pin or verification code")
	}

	currentTime := time.Now().UTC().Add(7 * time.Hour)
	expiredTimestamp := expiredCode.Unix()
	currentTimestamp := currentTime.Unix()
	if currentTimestamp > expiredTimestamp {
		fmt.Println("Verification code has expired")
		return resp, errors.New("verification code has expired")
	}

	return resp, nil
}

func (repo *userRepository) UserCreate(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, err error) {
	query := "INSERT INTO users (fullname, username, email, pin, phone_number) VALUES ($1, $2, $3, $4, $5) RETURNING id, fullname, username, email, phone_number"

	if err := repo.db.QueryRow(query, req.Fullname, req.Username, req.Email, req.Pin, req.PhoneNumber).Scan(&resp.Id, &resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber); err != nil {
		fmt.Println(err)
		return resp, errors.New("fail to create user")
	}

	return resp, nil
}

func (repo *userRepository) GetDataUserRepo(id string) (resp userDto.UserGetDataResponse, err error) {

	query := "SELECT fullname, username, email, phone_number FROM users WHERE id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber); err != nil {
		return resp, errors.New("fail to get data db")
	}

	return resp, nil
}

func (repo *userRepository) GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error) {

	query := "SELECT balance FROM users WHERE id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Balance); err != nil {
		return resp, errors.New("fail to get data db")
	}

	fmt.Println("Balance:", resp.Balance)

	return resp, nil
}

func (repo *userRepository) GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.TransactionRecord, error) {
	query := `
      SELECT
            t.id AS transaction_id,
            pm.payment_name,
            t.user_id,
            t.recipient_user_id,
            t.amount,
            t.description,
            t.transaction_date,
            t.status,
            u1.fullname AS sender_name,
            u2.fullname AS recipient_name
      FROM
            transaction t
      LEFT JOIN
            payment_method pm ON t.payment_method_id = pm.id
      LEFT JOIN
            users u1 ON t.user_id = u1.id
      LEFT JOIN
            users u2 ON t.recipient_user_id = u2.id
      WHERE
            (t.user_id = $1 OR t.recipient_user_id = $1)
   `

	var args []interface{}
	args = append(args, params.UserId)

	if params.TrxId != "" {
		query += " AND t.id = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.TrxId)
	}

	if params.TrxType != "" {
		if params.TrxType == "credit" {
			query += " AND (pm.payment_name IS NOT NULL OR t.recipient_user_id = $1)"
		} else if params.TrxType == "debit" {
			query += " AND t.user_id = $1 AND pm.payment_name IS NULL"
		}
	}

	if params.TrxDateStart != "" {
		query += " AND t.transaction_date >= $" + strconv.Itoa(len(args)+1)
		args = append(args, params.TrxDateStart)
	}

	if params.TrxDateEnd != "" {
		query += " AND t.transaction_date <= $" + strconv.Itoa(len(args)+1)
		args = append(args, params.TrxDateEnd)
	}

	if params.TrxStatus != "" {
		query += " AND t.status = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.TrxStatus)
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []userDto.TransactionRecord
	for rows.Next() {
		var transaction userDto.TransactionRecord

		if err := rows.Scan(
			&transaction.TransactionId,
			&transaction.PaymentMethod,
			&transaction.UserId,
			&transaction.RecipientUserId,
			&transaction.Amount,
			&transaction.Description,
			&transaction.TransactionDate,
			&transaction.PaymentStatus,
			&transaction.SenderName,
			&transaction.RecipientName,
		); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found for the given user ID or recipient user ID")
	}

	return transactions, nil
}
