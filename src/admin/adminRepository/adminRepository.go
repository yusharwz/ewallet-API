package adminRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/model/dto/adminDto"
	"time"

	"fmt"
	"strconv"
)

type adminRepo struct {
	db *sql.DB
}

func (r *adminRepo) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	query := "SELECT id, fullname, username, email, phone_number, created_at FROM users WHERE 1=1 AND deleted_at IS NULL"
	var args []interface{}
	argIndex := 1

	if params.ID != "" {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, params.ID)
		argIndex++
	}

	if params.Fullname != "" {
		query += fmt.Sprintf(" AND fullname LIKE $%d", argIndex)
		args = append(args, "%"+params.Fullname+"%")
		argIndex++
	}

	if params.Email != "" {
		query += fmt.Sprintf(" AND email = $%d", argIndex)
		args = append(args, params.Email)
		argIndex++
	}

	if params.PhoneNumber != "" {
		query += fmt.Sprintf(" AND phone_number = $%d", argIndex)
		args = append(args, params.PhoneNumber)
		argIndex++
	}

	if params.CreateAt != "" {
		query += fmt.Sprintf(" AND created_at = $%d", argIndex)
		args = append(args, params.CreateAt)
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

	var users []adminDto.User
	for rows.Next() {
		var user adminDto.User
		err := rows.Scan(&user.ID, &user.Fullname, &user.Username, &user.Email, &user.PhoneNumber, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *adminRepo) SaveUser(user adminDto.User) error {
	// Check if username already exists
	var usernameExists bool
	usernameQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := r.db.QueryRow(usernameQuery, user.Username).Scan(&usernameExists)
	if err != nil {
		return err
	}
	if usernameExists {
		return errors.New("username already exists")
	}

	// Check if email already exists
	var emailExists bool
	emailQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err = r.db.QueryRow(emailQuery, user.Email).Scan(&emailExists)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("email already exists")
	}

	// Check if phone number already exists
	var phoneNumberExists bool
	phoneQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)"
	err = r.db.QueryRow(phoneQuery, user.PhoneNumber).Scan(&phoneNumberExists)
	if err != nil {
		return err
	}
	if phoneNumberExists {
		return errors.New("phone number already exists")
	}

	// Insert the user if all checks pass
	query := "INSERT INTO users(fullname, username, email, phone_number, pin, created_at) VALUES($1, $2, $3, $4, $5, $6)"
	_, err = r.db.Exec(query, user.Fullname, user.Username, user.Email, user.PhoneNumber, user.Pin, time.Now())
	if err != nil {
		return err
	}
	return nil
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
	// Check if the user ID is valid
	if user.ID == "" {
		return errors.New("invalid user ID")
	}

	// Check if the user exists
	var userExists bool
	userQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)"
	err := r.db.QueryRow(userQuery, user.ID).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !userExists {
		return errors.New("user does not exist")
	}

	// Check if email already exists for another user
	var emailExists bool
	emailQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id <> $2 AND deleted_at IS NULL)"
	err = r.db.QueryRow(emailQuery, user.Email, user.ID).Scan(&emailExists)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return errors.New("email already exists for another user")
	}

	// Check if phone number already exists for another user
	var phoneNumberExists bool
	phoneQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1 AND id <> $2 AND deleted_at IS NULL)"
	err = r.db.QueryRow(phoneQuery, user.PhoneNumber, user.ID).Scan(&phoneNumberExists)
	if err != nil {
		return fmt.Errorf("failed to check phone number existence: %w", err)
	}
	if phoneNumberExists {
		return errors.New("phone number already exists for another user")
	}

	// Update the user
	query := `
        UPDATE users 
        SET fullname = $1, username = $2, email = $3, phone_number = $4, pin = $5, updated_at = $6 
        WHERE id = $7 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, user.Fullname, user.Username, user.Email, user.PhoneNumber, user.Pin, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *adminRepo) GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error) {
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
	query := "SELECT id,user_id, balance,created_at FROM wallets WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if params.ID != "" {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, params.ID)
		argIndex++
	}

	if params.User_id != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, params.User_id)
		argIndex++
	}
	if params.MinBalance != nil {
		query += fmt.Sprintf(" AND balance >= $%d", argIndex)
		args = append(args, params.MinBalance)
		argIndex++
	}

	if params.MaxBalance != nil {
		query += fmt.Sprintf(" AND balance <= $%d", argIndex)
		args = append(args, params.MaxBalance)
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

	var wallets []adminDto.Wallet
	for rows.Next() {
		var wallet adminDto.Wallet
		err := rows.Scan(&wallet.ID, &wallet.User_id, &wallet.Balance, &wallet.CreatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}
func NewAdminRepository(db *sql.DB) *adminRepo {
	return &adminRepo{db}
}
