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
			log.Error().Msg("invalid start date format: %s"+ err.Error())
			return nil, fmt.Errorf("invalid start date format: %s", err.Error())
		}
		endDate, err := time.Parse("2006-01-02", params.EndDate)
		if err != nil {
			log.Error().Msg("invalid start date format: %s"+ err.Error())
			return nil, fmt.Errorf("invalid end date format: %s", err.Error())
		}
		query += " AND created_at BETWEEN $" + strconv.Itoa(len(args)+1) + " AND $" + strconv.Itoa(len(args)+2)
		args = append(args, startDate, endDate)
	}

	if params.Page != "" && params.Limit != "" {
		page, err := strconv.Atoi(params.Page)
		if err != nil {
			log.Error().Msg("invalid page parameter: %s"+ err.Error())
			return nil, fmt.Errorf("invalid page parameter: %s", err.Error())
		}
		limit, err := strconv.Atoi(params.Limit)
		if err != nil {
			log.Error().Msg("invalid limit parameter: %s"+ err.Error())
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

	var users []adminDto.User
	for rows.Next() {
		var user adminDto.User
		err := rows.Scan(&user.ID, &user.Fullname, &user.Username, &user.ImageURL, &user.Pin, &user.Email, &user.PhoneNumber, &user.Roles, &user.Status, &user.CreatedAt)
		if err != nil {
			return nil, err
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

func (r *adminRepo) GetTransaction(page int, limit int) ([]adminDto.Transaction, int, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query("select id, user_id, transaction_type, amount, description, status, created_at from transactions LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data :%w", err)
	}
	defer rows.Close()

	var transactions []adminDto.Transaction
	for rows.Next() {
		var data adminDto.Transaction
		err := rows.Scan(&data.Id, &data.UserId, &data.TransactionType, &data.Amount, &data.Description, &data.Status, &data.Created_at)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan transaction data : %w", err)
		}

		var transDetails []adminDto.TransactionDetail
		rowDetails, err := r.db.Query("select id, transaction_id, wallet_transaction_id,topup_transaction_id, created_at from transactions_detail", data.Id)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to querytransaction detail data: %w", err)
		}
		defer rowDetails.Close()

		for rowDetails.Next() {
			var transDetail adminDto.TransactionDetail
			var walletTransactionID, topupTransactionID sql.NullString
			err := rowDetails.Scan(&transDetail.Id, &transDetail.TransactionId, &walletTransactionID, &topupTransactionID, &transDetail.Created_at)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to scan transaction detail data: %w", err)
			}

			if walletTransactionID.Valid {
				transDetail.WalletTransactionId = walletTransactionID.String
			}
			if topupTransactionID.Valid {
				transDetail.TopUpTransactionId = topupTransactionID.String
			}

			transDetails = append(transDetails, transDetail)
		}

		data.TransactionDetail = transDetails

		transactions = append(transactions, data)
	}

	total := len(transactions)

	return transactions, total, nil
}

func (r *adminRepo) GetWalletTransaction(page int, limit int) ([]adminDto.WalletTransaction, int, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query("select id, transaction_id, from_wallet_id, to_wallet_id, created_at from wallet_transactions limit $1 offset $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query wallet transaction data : %w", err)
	}
	defer rows.Close()

	var transactionsWallet []adminDto.WalletTransaction
	for rows.Next() {
		var data adminDto.WalletTransaction
		err := rows.Scan(&data.Id, &data.TransactionId, &data.FromWalletId, &data.ToWalletId, &data.Created_at)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan wallet transanction data:%w", err)
		}
		transaction := adminDto.WalletTransaction{
			Id:            data.Id,
			TransactionId: data.TransactionId,
			FromWalletId:  data.FromWalletId,
			ToWalletId:    data.ToWalletId,
			Created_at:    data.Created_at,
		}

		transactionsWallet = append(transactionsWallet, transaction)
	}

	total := len(transactionsWallet)
	return transactionsWallet, total, nil
}

func (r *adminRepo) GetTopUpTransaction(page int, limit int) ([]adminDto.TopUpTransaction, int, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query("select id, transaction_id, payment_method_id, created_at from topup_transactions limit $1 offset $2", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query get topup transaction data:%w", err)
	}
	defer rows.Close()

	var transactionsTopUp []adminDto.TopUpTransaction
	for rows.Next() {
		var data adminDto.TopUpTransaction
		err := rows.Scan(&data.Id, &data.TransactionId, &data.PaymentMethodId, &data.Created_at)
		if err != nil {
			return nil, 0, err
		}
		transaction := adminDto.TopUpTransaction{
			Id:              data.Id,
			TransactionId:   data.TransactionId,
			PaymentMethodId: data.PaymentMethodId,
			Created_at:      data.Created_at,
		}

		transactionsTopUp = append(transactionsTopUp, transaction)

	}

	total := len(transactionsTopUp)
	return transactionsTopUp, total, nil
}
