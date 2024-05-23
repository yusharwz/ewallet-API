package userRepository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/src/user"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-resty/resty/v2"
)

type userRepository struct {
	db     *sql.DB
	client *resty.Client
}

func NewUserRepository(db *sql.DB, client *resty.Client) user.UserRepository {
	return &userRepository{
		db:     db,
		client: client,
	}
}

func (repo *userRepository) CekEmail(email string) (bool, error) {
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

func (repo *userRepository) CekPhoneNumber(pnumber string) (bool, error) {
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

func (repo *userRepository) InsertCode(code, email, pnumber string) (bool, error) {
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

func (repo *userRepository) UserLogin(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {
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

func (repo *userRepository) UserUploadImage(req userDto.UploadImagesRequest) (userDto.UploadImagesResponse, error) {

	cldService, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	ctx := context.Background()

	var resp userDto.UploadImagesResponse
	response, _ := cldService.Upload.Upload(ctx, req.File, uploader.UploadParams{})

	resp.Url = response.SecureURL

	return resp, nil
}

func (repo *userRepository) ImageToDB(userId string, req userDto.UploadImagesResponse) error {
	query := `
		UPDATE users
		SET image_url = $1
		WHERE id = $2
	`
	if _, err := repo.db.Exec(query, req.Url, userId); err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) UserCreate(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, unique string, err error) {

	checkEmailQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
	var emailCount int
	if err := repo.db.QueryRow(checkEmailQuery, req.Email).Scan(&emailCount); err != nil {
		return resp, "", errors.New("failed to check email")
	}
	if emailCount > 0 {
		return resp, "", errors.New("email is already in use")
	}

	checkUsernameQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
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

func (repo *userRepository) ActivedAccount(req userDto.ActivatedAccountReq) (err error) {

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

func (repo *userRepository) UserWalletCreate(id string) (err error) {
	query := "INSERT INTO wallets (user_id) VALUES ($1)"

	if _, err := repo.db.Exec(query, id); err != nil {
		return errors.New("fail to create wallet")
	}

	return nil
}

func (repo *userRepository) GetDataUserRepo(id string) (resp userDto.UserGetDataResponse, err error) {
	var images sql.NullString
	query := "SELECT fullname, username, email, phone_number, image_url FROM users WHERE id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber, &images); err != nil {
		return resp, errors.New("fail to get data db")
	}

	if images.Valid {
		resp.ProfilImages = images.String
	}

	return resp, nil
}

func (repo *userRepository) GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error) {

	query := "SELECT balance FROM wallets WHERE user_id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Balance); err != nil {
		return resp, errors.New("fail to get data db")
	}

	return resp, nil
}

func (repo *userRepository) GetTotalDataCount(params userDto.GetTransactionParams) (totalData int, err error) {

	subquery1 := `
      SELECT COUNT(*)
      FROM
            transactions t
      WHERE
            t.user_id = $1
   `

	subquery2 := `
      SELECT COUNT(*)
      FROM
            transactions t
      JOIN
            wallet_transactions wt ON t.id = wt.transaction_id
      JOIN
            wallets w ON wt.from_wallet_id = w.id OR wt.to_wallet_id = w.id
      WHERE
            w.user_id = $1
   `

	args := []interface{}{params.UserId}
	conditionIndex := 2

	addCondition := func(condition string, value interface{}) {
		subquery1 += fmt.Sprintf(" AND %s = $%d", condition, conditionIndex)
		subquery2 += fmt.Sprintf(" AND %s = $%d", condition, conditionIndex)
		args = append(args, value)
		conditionIndex++
	}

	if params.TrxId != "" {
		addCondition("t.id", params.TrxId)
	}
	if params.TrxType != "" {
		addCondition("t.transaction_type", params.TrxType)
	}
	if params.TrxDateStart != "" {
		subquery1 += fmt.Sprintf(" AND t.created_at >= $%d", conditionIndex)
		subquery2 += fmt.Sprintf(" AND t.created_at >= $%d", conditionIndex)
		args = append(args, params.TrxDateStart)
		conditionIndex++
	}
	if params.TrxDateEnd != "" {
		subquery1 += fmt.Sprintf(" AND t.created_at <= $%d", conditionIndex)
		subquery2 += fmt.Sprintf(" AND t.created_at <= $%d", conditionIndex)
		args = append(args, params.TrxDateEnd)
		conditionIndex++
	}
	if params.TrxStatus != "" {
		addCondition("t.status", params.TrxStatus)
	}

	query := `
      SELECT SUM(count)
      FROM (
            ` + subquery1 + `
            UNION ALL
            ` + subquery2 + `
      ) AS subquery
   `

	if err := repo.db.QueryRow(query, args...).Scan(&totalData); err != nil {
		return 0, fmt.Errorf("fail to get total data count: %w", err)
	}

	return totalData, nil
}

func (repo *userRepository) GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, error) {

	baseQuery := `
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
	args := []interface{}{params.UserId}
	conditionIndex := 2

	addCondition := func(condition string, value interface{}) {
		baseQuery += fmt.Sprintf(" AND %s = $%d", condition, conditionIndex)
		args = append(args, value)
		conditionIndex++
	}

	if params.TrxId != "" {
		addCondition("t.id", params.TrxId)
	}
	if params.TrxType != "" {
		addCondition("t.transaction_type", params.TrxType)
	}
	if params.TrxDateStart != "" {
		baseQuery += fmt.Sprintf(" AND t.created_at >= $%d", conditionIndex)
		args = append(args, params.TrxDateStart)
		conditionIndex++
	}
	if params.TrxDateEnd != "" {
		baseQuery += fmt.Sprintf(" AND t.created_at <= $%d", conditionIndex)
		args = append(args, params.TrxDateEnd)
		conditionIndex++
	}
	if params.TrxStatus != "" {
		addCondition("t.status", params.TrxStatus)
	}

	if params.Page != "" && params.Limit != "" {
		page, _ := strconv.Atoi(params.Page)
		limit, _ := strconv.Atoi(params.Limit)
		offset := (page - 1) * limit
		baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	query := `
      SELECT
            id,
            transaction_type,
            amount,
            description,
            created_at,
            status
      FROM (
            ` + baseQuery + `
      ) sub
      WHERE 1=1
   `

	rows, err := repo.db.Query(query, args...)
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

		transaction.Detail = userDto.TransactionDetail{}

		paymentMethodQuery := `
            SELECT
               pm.payment_name, tt.payment_url
            FROM
               topup_transactions tt
            JOIN
               payment_method pm ON tt.payment_method_id = pm.id
            WHERE
               tt.transaction_id = $1
      `
		var paymentMethod sql.NullString
		var paymentURL sql.NullString
		err = repo.db.QueryRow(paymentMethodQuery, transaction.TransactionId).Scan(&paymentMethod, &paymentURL)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to query topup transaction: %w", err)
		}
		if paymentMethod.Valid {
			transaction.Detail.PaymentMethod = paymentMethod.String
			if paymentURL.Valid {
				transaction.Detail.PaymentURL = paymentURL.String
			}
		}

		walletTransactionQuery := `
            SELECT
               (SELECT u.username FROM users u JOIN wallets wf ON u.id = wf.user_id WHERE wf.id = wt.from_wallet_id) AS sender_name,
               (SELECT u.id FROM users u JOIN wallets wf ON u.id = wf.user_id WHERE wf.id = wt.from_wallet_id) AS sender_id,
               (SELECT u.username FROM users u JOIN wallets wr ON u.id = wr.user_id WHERE wr.id = wt.to_wallet_id) AS recipient_name,
               (SELECT u.id FROM users u JOIN wallets wr ON u.id = wr.user_id WHERE wr.id = wt.to_wallet_id) AS recipient_id
            FROM
               wallet_transactions wt
            WHERE
               wt.transaction_id = $1
      `
		var senderName, recipientName, senderId, recipientId sql.NullString
		err = repo.db.QueryRow(walletTransactionQuery, transaction.TransactionId).Scan(&senderName, &senderId, &recipientName, &recipientId)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to query wallet transaction: %w", err)
		}
		if senderName.Valid {
			transaction.Detail.SenderName = senderName.String
		}
		if recipientName.Valid {
			transaction.Detail.RecipientName = recipientName.String
		}
		if senderId.Valid {
			transaction.Detail.SenderId = senderId.String
		}
		if recipientId.Valid {
			transaction.Detail.RecipientId = recipientId.String
		}

		resp = append(resp, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over transaction rows: %w", err)
	}

	if len(resp) == 0 {
		return nil, errors.New("transaction not found")
	}

	return resp, nil
}

func (repo *userRepository) GetPaymentMethodName(id string) (metdhodName string, err error) {

	query := "SELECT payment_name FROM payment_method WHERE id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&metdhodName); err != nil {
		return "", errors.New("fail to get payment method name")
	}

	return metdhodName, nil
}

func (repo *userRepository) GetUserFullname(id string) (userFullname string, err error) {

	query := "SELECT fullname FROM users WHERE id = $1;"
	if err := repo.db.QueryRow(query, id).Scan(&userFullname); err != nil {
		return "", errors.New("fail to get user fullname")
	}

	return userFullname, nil
}

func (repo *userRepository) CreateTopUpTransaction(req userDto.TopUpTransactionRequest) (string, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return "", err
	}

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
		return "", err
	}

	transactionQuery := `
		INSERT INTO transactions (user_id, transaction_type, amount, description, created_at, status)
		VALUES ($1, 'credit', $2, $3, $4, 'pending')
		RETURNING id
	`
	var transactionID string
	err = tx.QueryRow(transactionQuery, req.UserId, req.Amount, req.Description, time.Now()).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	topupTransactionQuery := `
		INSERT INTO topup_transactions (transaction_id, payment_method_id, created_at)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(topupTransactionQuery, transactionID, req.PaymentMethodId, time.Now())
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return transactionID, nil
}

func (repo *userRepository) PaymentGateway(payload userDto.MidtransSnapReq) (userDto.MidtransSnapResp, error) {
	url := "https://app.sandbox.midtrans.com/snap/v1/transactions"
	serverKey := os.Getenv("SERVER_KEY")
	encodeKey := base64.StdEncoding.EncodeToString([]byte(serverKey))

	resp, err := repo.client.R().SetHeader("Authorization", "Basic "+encodeKey).SetBody(payload).Post(url)

	if err != nil {
		return userDto.MidtransSnapResp{}, err
	}

	var snapResp userDto.MidtransSnapResp
	err = json.Unmarshal(resp.Body(), &snapResp)
	if err != nil {
		return userDto.MidtransSnapResp{}, err
	}

	redirectUrl := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%s", snapResp.Token)
	snapResp.RedirectUrl = redirectUrl

	return snapResp, nil
}

func (repo *userRepository) InsertURL(transactionId, url string) (err error) {

	query := `
		UPDATE topup_transactions
		SET payment_url = $1
		WHERE transaction_id = $2
	`
	if _, err := repo.db.Exec(query, url, transactionId); err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) UpdateTransactionStatus(orderID string, status string) error {

	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := repo.db.Exec(query, status, orderID)
	if err != nil {
		return errors.New("failed to update transaction status")
	}

	return err
}

func (repo *userRepository) UpdateBalance(orderID, amountStr string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}

	var userID string
	query := `SELECT user_id FROM transactions WHERE id = $1`
	err = repo.db.QueryRow(query, orderID).Scan(&userID)
	if err != nil {
		return err
	}

	var walletID string
	query = `SELECT id FROM wallets WHERE user_id = $1`
	err = repo.db.QueryRow(query, userID).Scan(&walletID)
	if err != nil {
		return err
	}

	query = `UPDATE wallets SET balance = balance + $1 WHERE id = $2`
	_, err = repo.db.Exec(query, amount, walletID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) CreateWalletTransaction(req userDto.WalletTransactionRequest) (userDto.WalletTransactionResponse, string, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return userDto.WalletTransactionResponse{}, "", err
	}

	getWalletIdQuery := `SELECT id FROM wallets WHERE user_id = $1`
	err = tx.QueryRow(getWalletIdQuery, req.UserId).Scan(&req.FromWalletId)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", fmt.Errorf("disana: %v", err)
	}

	validatePinQuery := `SELECT pin FROM users WHERE id = $1`
	var storedPin string
	err = tx.QueryRow(validatePinQuery, req.UserId).Scan(&storedPin)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("user not found")
	}

	getRecipientWalletIdQuery := `
      SELECT w.id 
      FROM wallets w
      JOIN users u ON u.id = w.user_id
      WHERE u.phone_number = $1
   `
	err = tx.QueryRow(getRecipientWalletIdQuery, req.RecipientPhoneNumber).Scan(&req.ToWalletId)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("recipient not found")
	}

	var senderBalance float64
	balanceQuery := `SELECT balance FROM wallets WHERE id = $1`
	err = tx.QueryRow(balanceQuery, req.FromWalletId).Scan(&senderBalance)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", fmt.Errorf("disini: %v", err)
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return userDto.WalletTransactionResponse{}, "", fmt.Errorf("invalid amount: %v", err)
	}

	if senderBalance < amount {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("insufficient balance")
	}

	var recipientBalance float64
	err = tx.QueryRow(balanceQuery, req.ToWalletId).Scan(&recipientBalance)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("recipient wallet not found")
	}

	transactionQuery := `
      INSERT INTO transactions (user_id, transaction_type, amount, description, created_at, status)
      VALUES ($1, 'debit', $2, $3, $4, 'success')
      RETURNING id
   `
	var transactionID string
	err = tx.QueryRow(transactionQuery, req.UserId, req.Amount, req.Description, time.Now()).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", err
	}

	walletTransactionQuery := `
      INSERT INTO wallet_transactions (transaction_id, from_wallet_id, to_wallet_id, created_at)
      VALUES ($1, $2, $3, $4)
   `
	_, err = tx.Exec(walletTransactionQuery, transactionID, req.FromWalletId, req.ToWalletId, time.Now())
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", err
	}

	currentTime := time.Now()

	updateSenderBalanceQuery := `
      UPDATE wallets
      SET balance = balance - $1, updated_at = $2
      WHERE id = $3 AND balance >= $1
   `

	res, err := tx.Exec(updateSenderBalanceQuery, amount, currentTime, req.FromWalletId)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", err
	}

	if rowsAffected < 1 {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("insufficient balance after check")
	}

	updateRecipientBalanceQuery := `
      UPDATE wallets
      SET balance = balance + $1, updated_at = $2
      WHERE id = $3
   `
	_, err = tx.Exec(updateRecipientBalanceQuery, amount, currentTime, req.ToWalletId)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", err
	}

	if err := tx.Commit(); err != nil {
		return userDto.WalletTransactionResponse{}, "", err
	}

	return userDto.WalletTransactionResponse{TransactionId: transactionID}, storedPin, nil
}
