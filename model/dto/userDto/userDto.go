package userDto

import "database/sql"

type (
	UserLoginCodeRequestEmail struct {
		Email string `json:"email" binding:"required,email"`
	}

	UserLoginCodeRequestPhoneNumber struct {
		PhoneNumber string `json:"phoneNumber" binding:"required,nomorHp"`
	}

	UserLoginRequest struct {
		PhoneNumber string `json:"phoneNumber"`
		Pin         string `json:"pin" binding:"required"`
		Code        string `json:"code" binding:"required"`
	}

	UserLoginResponse struct {
		UserId    string `json:"userId,omitempty"`
		UserEmail string `json:"userEmail,omitempty"`
		Pin       string `json:"pin,omitempty"`
		Token     string `json:"token,omitempty"`
	}

	UserCreateRequest struct {
		Fullname    string `json:"fullname" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Pin         string `json:"pin" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required,nomorHp"`
	}

	UserCreateResponse struct {
		Id          string `json:"id"`
		Fullname    string `json:"fullname"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}

	UserGetDataResponse struct {
		Fullname    string `json:"fullname,omitempty"`
		Username    string `json:"username,omitempty"`
		Email       string `json:"email,omitempty"`
		PhoneNumber string `json:"phoneNumber,omitempty"`
		Balance     string `json:"balance,omitempty"`
	}

	GetTransactionParams struct {
		UserId          string
		RecipientUserId string
		TrxId           string
		TrxType         string
		TrxDateStart    string
		TrxDateEnd      string
		TrxStatus       string
		Page            string
		Limit           string
	}

	TransactionRecord struct {
		TransactionId   string
		PaymentMethod   sql.NullString
		UserId          sql.NullString
		RecipientUserId sql.NullString
		Amount          string
		Description     string
		TransactionDate string
		PaymentStatus   string
		SenderName      sql.NullString
		RecipientName   sql.NullString
	}

	GetTransactionResponse struct {
		TransactionId   string `json:"transactionId,omitempty"`
		PaymentMethod   string `json:"paymentMethod,omitempty"`
		TransactionType string `json:"transactionType,omitempty"`
		RecipientName   string `json:"recipientName,omitempty"`
		SenderName      string `json:"senderName,omitempty"`
		Amount          string `json:"amount,omitempty"`
		Description     string `json:"description"`
		TransactionDate string `json:"transactionDate"`
		PaymentStatus   string `json:"paymentStatus"`
	}
)
