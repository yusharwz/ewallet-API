package userRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/user"
	"fmt"
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

func (repo *userRepository) UserWalletCreate(id string) (err error) {
	query := "INSERT INTO wallets (user_id) VALUES ($1)"

	if _, err := repo.db.Exec(query, id); err != nil {
		fmt.Println(err)
		return errors.New("fail to create wallet")
	}

	return nil
}

func (repo *userRepository) GetDataUserRepo(id string) (resp userDto.UserGetDataResponse, err error) {

	query := "SELECT fullname, username, email, phone_number FROM users WHERE id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber); err != nil {
		return resp, errors.New("fail to get data db")
	}

	return resp, nil
}

func (repo *userRepository) GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error) {

	query := "SELECT balance FROM wallets WHERE user_id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Balance); err != nil {
		return resp, errors.New("fail to get data db")
	}

	fmt.Println("Balance:", resp.Balance)

	return resp, nil
}

func (repo *userRepository) GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, error) {
	// Query utama untuk transaksi pengguna
	query := `
      SELECT
            t.id,
            t.transaction_type,
            t.amount,
            t.description,
            t.created_at,
            t.status
      FROM
            transactions t
      WHERE
            t.user_id = $1
      UNION
      SELECT
            t.id,
            t.transaction_type,
            t.amount,
            t.description,
            t.created_at,
            t.status
      FROM
            transactions t
      JOIN
            wallet_transactions wt ON t.id = wt.transaction_id
      JOIN
            wallets w ON wt.from_wallet_id = w.id OR wt.to_wallet_id = w.id
      WHERE
            w.user_id = $1
   `

	rows, err := repo.db.Query(query, params.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from db: %w", err)
	}
	defer rows.Close()

	var resp []userDto.GetTransactionResponse

	for rows.Next() {
		var transaction userDto.GetTransactionResponse
		if err := rows.Scan(&transaction.TransactionId, &transaction.TransactionType, &transaction.Amount, &transaction.Description, &transaction.TransactionDate, &transaction.Status); err != nil {
			return nil, fmt.Errorf("failed to scan transaction data: %w", err)
		}

		// Default empty detail
		transaction.Detail = userDto.TransactionDetail{}

		// Query untuk mendapatkan payment_method_id dari topup_transactions
		paymentMethodQuery := `
            SELECT
               pm.payment_name
            FROM
               topup_transactions tt
            JOIN
               payment_method pm ON tt.payment_method_id = pm.id
            WHERE
               tt.transaction_id = $1
      `
		var paymentMethod sql.NullString
		err = repo.db.QueryRow(paymentMethodQuery, transaction.TransactionId).Scan(&paymentMethod)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to query topup transaction: %w", err)
		}
		if paymentMethod.Valid {
			transaction.Detail.PaymentMethod = paymentMethod.String
		}

		// Query untuk mendapatkan SenderName dan RecipientName dari wallet_transactions
		walletTransactionQuery := `
            SELECT
                (SELECT u.username FROM users u JOIN wallets wf ON u.id = wf.user_id WHERE wf.id = wt.from_wallet_id) AS sender_name,
                (SELECT u.username FROM users u JOIN wallets wr ON u.id = wr.user_id WHERE wr.id = wt.to_wallet_id) AS recipient_name
            FROM
                wallet_transactions wt
            WHERE
                wt.transaction_id = $1
        `
		var senderName, recipientName sql.NullString
		err = repo.db.QueryRow(walletTransactionQuery, transaction.TransactionId).Scan(&senderName, &recipientName)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to query wallet transaction: %w", err)
		}
		if senderName.Valid {
			transaction.Detail.SenderName = senderName.String
		}
		if recipientName.Valid {
			transaction.Detail.RecipientName = recipientName.String
		}

		resp = append(resp, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over transaction rows: %w", err)
	}

	return resp, nil
}

func (repo *userRepository) CreateTopUpTransaction(req userDto.TopUpTransactionRequest) (userDto.TopUpTransactionResponse, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return userDto.TopUpTransactionResponse{}, err
	}

	// Check if payment_method_id is valid
	var validPaymentMethod bool
	checkPaymentMethodQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM payment_method
			WHERE id = $1
		)
	`
	err = tx.QueryRow(checkPaymentMethodQuery, req.PaymentMethodId).Scan(&validPaymentMethod)
	if err != nil {
		tx.Rollback()
		return userDto.TopUpTransactionResponse{}, err
	}

	// Insert into transactions table without specifying the id
	transactionQuery := `
		INSERT INTO transactions (user_id, transaction_type, amount, description, created_at, status)
		VALUES ($1, 'credit', $2, $3, $4, 'pending')
		RETURNING id
	`
	var transactionID string
	err = tx.QueryRow(transactionQuery, req.UserId, req.Amount, req.Description, time.Now()).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return userDto.TopUpTransactionResponse{}, err
	}

	// Insert into topup_transactions table without specifying the id
	topupTransactionQuery := `
		INSERT INTO topup_transactions (transaction_id, payment_method_id, created_at)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(topupTransactionQuery, transactionID, req.PaymentMethodId, time.Now())
	if err != nil {
		tx.Rollback()
		return userDto.TopUpTransactionResponse{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return userDto.TopUpTransactionResponse{}, err
	}

	return userDto.TopUpTransactionResponse{TransactionId: transactionID}, nil
}

func (repo *userRepository) CreateWalletTransaction(req userDto.WalletTransactionRequest) (userDto.WalletTransactionResponse, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return userDto.WalletTransactionResponse{}, err
	}

	transactionQuery := `
		INSERT INTO transactions (user_id, transaction_type, amount, description, created_at, status)
		VALUES ($1, 'debit', $2, $3, $4, 'pending')
		RETURNING id
	`
	var transactionID string
	err = tx.QueryRow(transactionQuery, req.UserId, req.Amount, req.Description, time.Now()).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, err
	}

	// Insert into wallet_transactions table without specifying the id
	walletTransactionQuery := `
		INSERT INTO wallet_transactions (transaction_id, from_wallet_id, to_wallet_id, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(walletTransactionQuery, transactionID, req.FromWalletId, req.ToWalletId, time.Now())
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return userDto.WalletTransactionResponse{}, err
	}

	return userDto.WalletTransactionResponse{TransactionId: transactionID}, nil
}
