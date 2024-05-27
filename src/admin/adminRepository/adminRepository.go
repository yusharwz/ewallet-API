package adminRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/adminDto"
	"time"

	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *adminRepo {
	return &adminRepo{db}
}

func (r *adminRepo) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	query := "SELECT id, fullname, username, image_url, pin, email, phone_number, roles, status, created_at FROM users WHERE deleted_at IS NULL"
	args := []interface{}{}

	if params.ID != "" {
		query += " AND id = $1"
		args = append(args, params.ID)
	}

	if params.Fullname != "" {
		query += " AND fullname LIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+params.Fullname+"%")
	}

	if params.Username != "" {
		query += " AND username = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.Username)
	}

	if params.Email != "" {
		query += " AND email = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.Email)
	}

	if params.PhoneNumber != "" {
		query += " AND phone_number = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.PhoneNumber)
	}

	if params.Roles != "" {
		query += " AND roles = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.Roles)
	}

	if params.Status != "" {
		query += " AND status = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.Status)
	}

	if params.StartDate != "" && params.EndDate != "" {
		startDate, err := time.Parse("2006-01-02", params.StartDate)
		if err != nil {
			log.Error().Msg("invalid start date format: %s" + err.Error())
			return nil, fmt.Errorf("invalid start date format: %s", err.Error())
		}
		endDate, err := time.Parse("2006-01-02", params.EndDate)
		if err != nil {
			log.Error().Msg("invalid start date format: %s" + err.Error())
			return nil, fmt.Errorf("invalid end date format: %s", err.Error())
		}
		query += " AND created_at BETWEEN $" + strconv.Itoa(len(args)+1) + " AND $" + strconv.Itoa(len(args)+2)
		args = append(args, startDate, endDate)
	}

	if params.Page != "" && params.Limit != "" {
		page, err := strconv.Atoi(params.Page)
		if err != nil {
			log.Error().Msg("invalid page parameter: %s" + err.Error())
			return nil, fmt.Errorf("invalid page parameter: %s", err.Error())
		}
		limit, err := strconv.Atoi(params.Limit)
		if err != nil {
			log.Error().Msg("invalid limit parameter: %s" + err.Error())
			return nil, fmt.Errorf("invalid limit parameter: %s", err.Error())
		}
		offset := (page - 1) * limit
		query += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var image sql.NullString

	var users []adminDto.User
	for rows.Next() {
		var user adminDto.User
		err := rows.Scan(&user.ID, &user.Fullname, &user.Username, &image, &user.Pin, &user.Email, &user.PhoneNumber, &user.Roles, &user.Status, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		if image.Valid {
			user.ImageURL = image.String
		}
		users = append(users, user)
	}

	if params.Username != "" && len(users) == 0 {
		return nil, fmt.Errorf("user with username '%s' not found", params.Username)
	}
	if params.Fullname != "" && len(users) == 0 {
		return nil, fmt.Errorf("user with fullname '%s' not found", params.Fullname)
	}
	if params.Email != "" && len(users) == 0 {
		return nil, fmt.Errorf("user with email '%s' not found", params.Email)
	}

	if params.PhoneNumber != "" && len(users) == 0 {
		return nil, fmt.Errorf("user with phone number '%s' not found", params.PhoneNumber)
	}

	if params.ID != "" && len(users) == 0 {
		return nil, fmt.Errorf("user with ID '%s' not found", params.ID)
	}

	return users, nil
}

func (r *adminRepo) SoftDeleteUser(userID string) error {
	query := "UPDATE users SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *adminRepo) UpdateUser(user adminDto.User) error {
	if user.ID == "" {
		return errors.New("invalid user ID")
	}

	var userExists bool
	userQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)"
	err := r.db.QueryRow(userQuery, user.ID).Scan(&userExists)
	if err != nil {
		log.Error().Msg("failed to check user existence: %w" + err.Error())
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !userExists {
		log.Error().Msg("user does not exist")
		return errors.New("user does not exist")
	}
	var usernameExists bool
	usernameQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id <> $2 AND deleted_at IS NULL)"
	err = r.db.QueryRow(usernameQuery, user.Username, user.ID).Scan(&usernameExists)
	if err != nil {
		log.Error().Msg("failed to check username existence: %w" + err.Error())
		return fmt.Errorf("failed to check username existence: %w", err)
	}
	if usernameExists {
		log.Error().Msg("username already exists for another user")
		return errors.New("username already exists for another user")
	}
	var emailExists bool
	emailQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id <> $2 AND deleted_at IS NULL)"
	err = r.db.QueryRow(emailQuery, user.Email, user.ID).Scan(&emailExists)
	if err != nil {
		log.Error().Msg("failed to check email existence: %w" + err.Error())
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		log.Error().Msg("email already exists for another user")
		return errors.New("email already exists for another user")
	}
	var phoneNumberExists bool
	phoneQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1 AND id <> $2 AND deleted_at IS NULL)"
	err = r.db.QueryRow(phoneQuery, user.PhoneNumber, user.ID).Scan(&phoneNumberExists)
	if err != nil {
		log.Error().Msg("failed to check phone number existencer: %w" + err.Error())
		return fmt.Errorf("failed to check phone number existence: %w", err)
	}
	if phoneNumberExists {
		log.Error().Msg("phone number already exists for another user")
		return errors.New("phone number already exists for another user")
	}
	query := `
        UPDATE users
        SET fullname = $1, username = $2, email = $3, phone_number = $4, pin = $5, updated_at = $6
        WHERE id = $7 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, user.Fullname, user.Username, user.Email, user.PhoneNumber, user.Pin, time.Now(), user.ID)
	if err != nil {
		log.Error().Msg("failed to update user: %w" + err.Error())
		return fmt.Errorf("failed to update user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Msg("failed to get rows affected:%w" + err.Error())
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *adminRepo) GetpaymentMethodByParams(params adminDto.GetPaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	query := "SELECT id, payment_name,created_at FROM payment_method WHERE 1=1 AND deleted_at IS NULL"
	var args []interface{}
	argIndex := 1

	if params.ID != "" {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, params.ID)
		argIndex++
	}

	if params.PaymentName != "" {
		query += fmt.Sprintf(" AND payment_name LIKE $%d", argIndex)
		args = append(args, "%"+params.PaymentName+"%")
		argIndex++
	}

	if params.CreatedAt != "" {
		query += fmt.Sprintf(" AND created_at = $%d", argIndex)
		args = append(args, params.CreatedAt)
		argIndex++
	}
	if params.Page != "" && params.Limit != "" {
		page, err := strconv.Atoi(params.Page)
		if err != nil {
			return nil, fmt.Errorf("invalid page parameter")
		}
		limit, err := strconv.Atoi(params.Limit)
		if err != nil {
			return nil, fmt.Errorf("invalid limit parameter")
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, limit, offset)
		argIndex += 2
	}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentMethods []adminDto.PaymentMethod
	for rows.Next() {
		var paymentMethod adminDto.PaymentMethod
		err := rows.Scan(&paymentMethod.ID, &paymentMethod.PaymentName, &paymentMethod.CreatedAt)
		if err != nil {
			return nil, err
		}
		paymentMethods = append(paymentMethods, paymentMethod)
	}
	if params.ID != "" && len(paymentMethods) == 0 {
		return nil, fmt.Errorf("payment with id '%s' not found", params.ID)
	}

	if params.PaymentName != "" && len(paymentMethods) == 0 {
		return nil, fmt.Errorf("payment with name '%s' not found", params.PaymentName)
	}

	return paymentMethods, nil
}
func (r *adminRepo) checkPaymentMethodExists(paymentName string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM payment_method WHERE LOWER(payment_name) = LOWER($1) AND deleted_at IS NULL)"
	err := r.db.QueryRow(query, paymentName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *adminRepo) SavePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	exists, err := r.checkPaymentMethodExists(paymentMethod.PaymentName)
	if err != nil {
		return err
	}
	if exists {
		log.Error().Msg("payment method name already exists")
		return errors.New("payment method name already exists")
	}

	query := "INSERT INTO payment_method(payment_name) VALUES($1)"
	_, err = r.db.Exec(query, paymentMethod.PaymentName)
	if err != nil {
		return err
	}
	return nil
}
func (r *adminRepo) SoftDeletePaymentMethod(paymentMethodID string) error {
	query := "UPDATE payment_method SET deleted_at=$1 WHERE id=$2 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, time.Now(), paymentMethodID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *adminRepo) UpdatePaymentMethod(paymentMethod adminDto.PaymentMethod) error {
	exists, err := r.checkPaymentMethodExists(paymentMethod.PaymentName)
	if err != nil {
		return err
	}
	if exists {
		log.Error().Msg("payment method name already exists")
		return errors.New("payment method name already exists")
	}
	query := "UPDATE payment_method SET payment_name=$1, updated_at=$2 WHERE id=$3 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, paymentMethod.PaymentName, time.Now(), paymentMethod.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *adminRepo) GetWalletByParams(params adminDto.GetWalletParams) ([]adminDto.Wallet, error) {
	query := `
	SELECT w.id, w.user_id, w.balance, w.created_at, u.fullname, u.username
	FROM wallets w
	JOIN users u ON w.user_id = u.id
	WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if params.ID != "" {
		query += fmt.Sprintf(" AND w.id = $%d", argIndex)
		args = append(args, params.ID)
		argIndex++
	}

	if params.User_id != "" {
		query += fmt.Sprintf(" AND w.user_id = $%d", argIndex)
		args = append(args, params.User_id)
		argIndex++
	}
	if params.Fullname != "" {
		query += fmt.Sprintf(" AND u.fullname ILIKE $%d", argIndex)
		args = append(args, "%"+params.Fullname+"%")
		argIndex++
	}
	if params.Username != "" {
		query += fmt.Sprintf(" AND u.username ILIKE $%d", argIndex)
		args = append(args, "%"+params.Username+"%")
		argIndex++
	}
	if params.MinBalance != nil {
		query += fmt.Sprintf(" AND w.balance >= $%d", argIndex)
		args = append(args, params.MinBalance)
		argIndex++
	}

	if params.MaxBalance != nil {
		query += fmt.Sprintf(" AND w.balance <= $%d", argIndex)
		args = append(args, params.MaxBalance)
		argIndex++
	}

	if params.CreatedAt != "" {
		query += fmt.Sprintf(" AND w.created_at = $%d", argIndex)
		args = append(args, params.CreatedAt)
		argIndex++
	}
	if params.Page != "" && params.Limit != "" {
		page, err := strconv.Atoi(params.Page)
		if err != nil {
			log.Error().Msg("invalid page parameter")
			return nil, fmt.Errorf("invalid page parameter")
		}
		limit, err := strconv.Atoi(params.Limit)
		if err != nil {
			log.Error().Msg("invalid limt parameter")
			return nil, fmt.Errorf("invalid limit parameter")
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, limit, offset)
		argIndex += 2
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []adminDto.Wallet
	for rows.Next() {
		var wallet adminDto.Wallet
		err := rows.Scan(&wallet.ID, &wallet.User_id, &wallet.Balance, &wallet.CreatedAt, &wallet.Fullname, &wallet.Username)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	if params.Username != "" && len(wallets) == 0 {
		return nil, fmt.Errorf("wallets with username '%s' not found", params.Username)
	}

	if params.Fullname != "" && len(wallets) == 0 {
		return nil, fmt.Errorf("wallets with email '%s' not found", params.Fullname)
	}

	if params.ID != "" && len(wallets) == 0 {
		return nil, fmt.Errorf("wallets with ID '%s' not found", params.ID)
	}
	if params.User_id != "" && len(wallets) == 0 {
		return nil, fmt.Errorf("walets with user ID '%s' not found", params.User_id)
	}
	return wallets, nil
}

func (r *adminRepo) GetTransactionRepo(params adminDto.GetTransactionParams) ([]adminDto.GetTransactionResponse, int, error) {
	baseQuery := `
        SELECT
            t.id,
            t.amount,
            t.description,
            t.created_at,
            t.status,
            t.transaction_type,
            t.user_id,
            u.username
        FROM
            transactions t
        JOIN
            users u ON t.user_id = u.id
        WHERE
            1=1
    `

	args := []interface{}{}
	conditionIndex := 1

	addCondition := func(baseQuery *string, condition string, value interface{}) {
		*baseQuery += fmt.Sprintf(" AND %s $%d", condition, conditionIndex)
		args = append(args, value)
		conditionIndex++
	}

	if params.UserId != "" {
		addCondition(&baseQuery, "t.user_id =", params.UserId)
	}
	if params.TrxId != "" {
		addCondition(&baseQuery, "t.id =", params.TrxId)
	}
	if params.TrxDateStart != "" {
		addCondition(&baseQuery, "t.created_at >=", params.TrxDateStart)
	}
	if params.TrxDateEnd != "" {
		trxDateEnd := params.TrxDateEnd + " 23:59:59.999999"
		addCondition(&baseQuery, "t.created_at <=", trxDateEnd)
	}
	if params.TrxStatus != "" {
		addCondition(&baseQuery, "t.status ILIKE", "%"+params.TrxStatus+"%")
	}
	if params.TrxType != "" {
		addCondition(&baseQuery, "t.transaction_type =", params.TrxType)
	}

	finalQuery := baseQuery

	countQuery := `
        SELECT COUNT(*)
        FROM (
            ` + baseQuery + `
        ) sub
    `

	if params.Page != "" && params.Limit != "" {
		page, _ := strconv.Atoi(params.Page)
		limit, _ := strconv.Atoi(params.Limit)
		offset := (page - 1) * limit
		finalQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	var totalData int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total data count: %w", err)
	}

	rows, err := r.db.Query(finalQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data from db: %w", err)
	}
	defer rows.Close()

	var resp []adminDto.GetTransactionResponse

	for rows.Next() {
		var transaction adminDto.GetTransactionResponse
		if err := rows.Scan(&transaction.TransactionId, &transaction.Amount, &transaction.Description, &transaction.TransactionDate, &transaction.Status, &transaction.TransactionType, &transaction.UserId, &transaction.UserName); err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction data: %w", err)
		}

		transaction.Detail = adminDto.TransactionDetail{}

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
		err = r.db.QueryRow(paymentMethodQuery, transaction.TransactionId).Scan(&paymentMethod, &paymentURL)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, fmt.Errorf("failed to query topup transaction: %w", err)
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
		err = r.db.QueryRow(walletTransactionQuery, transaction.TransactionId).Scan(&senderName, &senderId, &recipientName, &recipientId)
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
		err = r.db.QueryRow(merchantTransactionQuery, transaction.TransactionId).Scan(&merchantName)
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
