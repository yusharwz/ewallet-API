package adminRepository

import (
	"database/sql"
	"final-project-enigma/model/dto/adminDto"

	"fmt"
	"strconv"
)

type adminRepo struct {
	db *sql.DB
}

func (r *adminRepo) GetUsersByParams(params adminDto.GetUserParams) ([]adminDto.User, error) {
	query := "SELECT id, fullname, username, email, phone_number, created_at FROM users WHERE 1=1"
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
func (r *adminRepo) GetpaymentMethodByParams(params adminDto.GetpaymentMethodParams) ([]adminDto.PaymentMethod, error) {
	query := "SELECT id, payment_name,created_at FROM payment_method WHERE 1=1"
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
