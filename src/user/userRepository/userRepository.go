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

func (repo *userRepository) EditUserData(req userDto.UserUpdateReq) error {
	usernameCheckQuery := `
		SELECT COUNT(*) 
		FROM users 
		WHERE username = $1 AND id != $2
	`
	var usernameCount int
	if err := repo.db.QueryRow(usernameCheckQuery, req.Username, req.UserId).Scan(&usernameCount); err != nil {
		return errors.New("failed to check username")
	}
	if usernameCount > 0 {
		return errors.New("username is already in use")
	}

	emailCheckQuery := `
		SELECT COUNT(*) 
		FROM users 
		WHERE email = $1 AND id != $2
	`
	var emailCount int
	if err := repo.db.QueryRow(emailCheckQuery, req.Email, req.UserId).Scan(&emailCount); err != nil {
		return errors.New("failed to check email")
	}
	if emailCount > 0 {
		return errors.New("email is already in use")
	}

	phoneCheckQuery := `
		SELECT COUNT(*) 
		FROM users 
		WHERE phone_number = $1 AND id != $2
	`
	var phoneCount int
	if err := repo.db.QueryRow(phoneCheckQuery, req.PhoneNumber, req.UserId).Scan(&phoneCount); err != nil {
		return errors.New("failed to check phone number")
	}
	if phoneCount > 0 {
		return errors.New("phone number is already in use")
	}

	updateQuery := `
		UPDATE users
		SET fullname = $1, username = $2, email = $3, phone_number = $4
		WHERE id = $5
	`
	if _, err := repo.db.Exec(updateQuery, req.Fullname, req.Username, req.Email, req.PhoneNumber, req.UserId); err != nil {
		return errors.New("failed to update user")
	}

	return nil
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

	currentTime := time.Now()

	query := `
		UPDATE users
		SET image_url = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`
	if _, err := repo.db.Exec(query, req.Url, currentTime, userId); err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) GetDataUserRepo(id string) (resp userDto.UserGetDataResponse, err error) {
	var images sql.NullString
	query := "SELECT fullname, username, email, phone_number, image_url FROM users WHERE id = $1 AND deleted_at IS NULL;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Fullname, &resp.Username, &resp.Email, &resp.PhoneNumber, &images); err != nil {
		return resp, errors.New("fail to get data db")
	}

	if images.Valid {
		resp.ProfilImages = images.String
	}

	return resp, nil
}

func (repo *userRepository) GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error) {

	query := "SELECT balance FROM wallets WHERE user_id = $1 AND deleted_at IS NULL;"
	if err := repo.db.QueryRow(query, id).Scan(&resp.Balance); err != nil {
		return resp, errors.New("fail to get data db")
	}

	return resp, nil
}

func (repo *userRepository) GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, int, error) {
	baseQuery1 := `
		SELECT
			t.id,
			t.amount,
			t.description,
			t.created_at,
			t.status
		FROM
			transactions t
		WHERE
			t.user_id = $1
	`

	baseQuery2 := `
		SELECT
			t.id,
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

	addCondition := func(baseQuery *string, condition string, value interface{}) {
		*baseQuery += fmt.Sprintf(" AND %s $%d", condition, conditionIndex)
		args = append(args, value)
		conditionIndex++
	}

	if params.TrxId != "" {
		addCondition(&baseQuery1, "t.id =", params.TrxId)
		addCondition(&baseQuery2, "t.id =", params.TrxId)
	}
	if params.TrxDateStart != "" {
		addCondition(&baseQuery1, "t.created_at >=", params.TrxDateStart)
		addCondition(&baseQuery2, "t.created_at >=", params.TrxDateStart)
	}
	if params.TrxDateEnd != "" {
		trxDateEnd := params.TrxDateEnd + " 23:59:59.999999"
		addCondition(&baseQuery1, "t.created_at <=", trxDateEnd)
		addCondition(&baseQuery2, "t.created_at <=", trxDateEnd)
	}
	if params.TrxStatus != "" {
		addCondition(&baseQuery1, "t.status ILIKE", "%"+params.TrxStatus+"%")
		addCondition(&baseQuery2, "t.status ILIKE", "%"+params.TrxStatus+"%")
	}

	finalQuery := `
		SELECT
			id,
			amount,
			description,
			created_at,
			status
		FROM (
			` + baseQuery1 + `
			UNION
			` + baseQuery2 + `
		) sub
		WHERE 1=1
	`

	countQuery := `
		SELECT COUNT(*)
		FROM (
			` + baseQuery1 + `
			UNION
			` + baseQuery2 + `
		) sub
		WHERE 1=1
	`

	if params.Page != "" && params.Limit != "" {
		page, _ := strconv.Atoi(params.Page)
		limit, _ := strconv.Atoi(params.Limit)
		offset := (page - 1) * limit
		finalQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	var totalData int
	err := repo.db.QueryRow(countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total data count: %w", err)
	}

	rows, err := repo.db.Query(finalQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from db: %w", err)
	}
	defer rows.Close()

	var resp []userDto.GetTransactionResponse

	for rows.Next() {
		var transaction userDto.GetTransactionResponse
		if err := rows.Scan(&transaction.TransactionId, &transaction.Amount, &transaction.Description, &transaction.TransactionDate, &transaction.Status); err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction data: %w", err)
		}

		transaction.Detail = userDto.TransactionDetail{}

		// Query for payment method details
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
			return nil, 0, fmt.Errorf("failed to query topup transaction: %w", err)
		}
		if paymentMethod.Valid {
			transaction.Detail.PaymentMethod = paymentMethod.String
			if paymentURL.Valid {
				transaction.Detail.PaymentURL = paymentURL.String
			}
		}

		// Query for wallet transaction details
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
			return nil, 0, fmt.Errorf("failed to query wallet transaction: %w", err)
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

		// Query for merchant transaction details
		merchantTransactionQuery := `
			SELECT
				m.merchant_name
			FROM
				merchant_transactions mt
			JOIN
				merchant m ON mt.merchant_id = m.id
			WHERE
				mt.transaction_id = $1
		`
		var merchantName sql.NullString
		err = repo.db.QueryRow(merchantTransactionQuery, transaction.TransactionId).Scan(&merchantName)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, fmt.Errorf("failed to query merchant transaction: %w", err)
		}
		if merchantName.Valid {
			transaction.Detail.MerchantName = merchantName.String
		}

		resp = append(resp, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate over transaction rows: %w", err)
	}

	if len(resp) == 0 {
		return nil, 0, errors.New("transaction not found")
	}

	return resp, totalData, nil
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
		return "", errors.New("payment method not registered")
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

func (repo *userRepository) InsertPaymentURL(transactionId, url string) (err error) {

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

	validatePinQuery := `SELECT pin FROM users WHERE id = $1 AND deleted_at IS NULL`
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
		WHERE u.phone_number = $1 AND u.status = 'active' AND u.deleted_at IS NULL
	`
	err = tx.QueryRow(getRecipientWalletIdQuery, req.RecipientPhoneNumber).Scan(&req.ToWalletId)
	if err != nil {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("recipient not found")
	}

	if req.FromWalletId == req.ToWalletId {
		tx.Rollback()
		return userDto.WalletTransactionResponse{}, "", errors.New("sender and recipient cannot be the same")
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

func (repo *userRepository) CreateMerchantTransaction(req userDto.MerchantTransactionRequest) (string, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return "", err
	}

	var validMerchant bool
	checkMerchantQuery := `
        SELECT EXISTS (
            SELECT 1
            FROM merchant
            WHERE id = $1
        )
    `
	err = tx.QueryRow(checkMerchantQuery, req.MerchantId).Scan(&validMerchant)
	if err != nil {
		tx.Rollback()
		return "", errors.New("payment method not registered")
	}

	if !validMerchant {
		tx.Rollback()
		return "", errors.New("invalid merchant")
	}

	var currentBalance float64
	checkBalanceQuery := `
        SELECT balance
        FROM wallets
        WHERE user_id = $1
    `
	err = tx.QueryRow(checkBalanceQuery, req.UserId).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if currentBalance < req.Amount {
		tx.Rollback()
		return "", errors.New("insufficient balance")
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
		return "", err
	}

	merchantTransactionQuery := `
        INSERT INTO merchant_transactions (transaction_id, merchant_id, created_at)
        VALUES ($1, $2, $3)
    `
	_, err = tx.Exec(merchantTransactionQuery, transactionID, req.MerchantId, time.Now())
	if err != nil {
		tx.Rollback()
		return "", err
	}

	updateBalanceQuery := `
        UPDATE wallets
        SET balance = balance - $1
        WHERE user_id = $2
    `
	_, err = tx.Exec(updateBalanceQuery, req.Amount, req.UserId)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return transactionID, nil
}

func (repo *userRepository) DeleteUser(id string) error {
	var deletedAt sql.NullTime

	checkStatusQuery := `SELECT deleted_at FROM users WHERE id = $1;`
	if err := repo.db.QueryRow(checkStatusQuery, id).Scan(&deletedAt); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	if deletedAt.Valid {
		return errors.New("user already deleted")
	}

	currentTime := time.Now()

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	// Update users table
	queryUsers := `
        UPDATE users
        SET deleted_at = $1, updated_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `
	if _, err := tx.Exec(queryUsers, currentTime, id); err != nil {
		tx.Rollback()
		return err
	}

	// Update wallets table
	queryWallets := `
        UPDATE wallets
        SET deleted_at = $1, updated_at = $1
        WHERE user_id = $2 AND deleted_at IS NULL
    `
	if _, err := tx.Exec(queryWallets, currentTime, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
