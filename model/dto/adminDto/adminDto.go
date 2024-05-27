package adminDto

import "time"

type (
	GetUserParams struct {
		ID          string `json:"id"`
		Fullname    string `json:"fullname"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		Roles       string `json:"roles"`
		Status      string `json:"status"`
		StartDate   string `json:"satrtDate"`
		EndDate     string `json:"endDate"`
		Page        string `json:"page"`
		Limit       string `json:"limit"`
	}
	User struct {
		ID          string    `json:"id"`
		Fullname    string    `json:"fullname"`
		Username    string    `json:"username"`
		Pin         string    `json:"pin"`
		Email       string    `json:"email"`
		ImageURL    string    `json:"image_url"`
		Roles       string    `json:"roles"`
		Status      string    `json:"status"`
		PhoneNumber string    `json:"phoneNumber"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}
	UserUpdateRequest struct {
		ID          string `json:"id"`
		Fullname    string `json:"fullname" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Pin         string `json:"pin" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	GetPaymentMethodParams struct {
		ID          string `json:"id"`
		PaymentName string `json:"payment_name"`
		CreatedAt   string `json:"createdAt"`
		Page        string `json:"page"`
		Limit       string `json:"limit"`
	}
	PaymentMethod struct {
		ID          string    `json:"id"`
		PaymentName string    `json:"payment_name"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	GetWalletParams struct {
		ID         string   `json:"id"`
		User_id    string   `json:"user_id"`
		Fullname   string   `json:"fullname"`
		Username   string   `json:"username"`
		MinBalance *float64 `json:"min_balance"`
		MaxBalance *float64 `json:"max_balance"`
		CreatedAt  string   `json:"createdAt"`
		Page       string   `json:"page"`
		Limit      string   `json:"limit"`
	}
	Wallet struct {
		ID       string `json:"id"`
		User_id  string `json:"user_id"`
		Fullname string `json:"fullname"`
		Username string `json:"username"`

		Balance   string    `json:"balance"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	CreatePaymentMethod struct {
		PaymentName string `json:"payment_name" binding:"required,max=255"`
	}

	UpdatePaymentRequest struct {
		ID          string `json:"id"`
		PaymentName string `json:"payment_name" binding:"required,max=255"`
	}

	PaymentResponse struct {
		ID          string `json:"id"`
		PaymentName string `json:"payment_name" binding:"required,max=255"`
	}

	PaymentMethodAdmin struct {
		Id          string `json:"id"`
		PaymentName string `json:"payment_name"`
	}

	TopUpTransaction struct {
		Id              string    `json:"id"`
		TransactionId   string    `json:"transaction_id"`
		PaymentMethodId string    `json:"payment_method_id"`
		Created_at      time.Time `json:"created_at"`
	}

	Transaction struct {
		Id                string              `json:"id"`
		UserId            string              `json:"user_id"`
		TransactionType   string              `json:"transaction_type"`
		Amount            float64             `json:"amount"`
		Description       string              `json:"description"`
		Status            string              `json:"status"`
		Created_at        time.Time           `json:"created_at"`
		TransactionDetail []TransactionDetail `json:"transactions_detail"`
	}

	WalletTransaction struct {
		Id            string    `json:"id"`
		TransactionId string    `json:"transaction_id"`
		FromWalletId  string    `json:"from_wallet_id"`
		ToWalletId    string    `json:"to_wallet_id"`
		Created_at    time.Time `json:"created_at"`
	}

	UserResponse struct {
		ID          string `json:"id"`
		Fullname    string `json:"fullname"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Roles       string `json:"roles"`
		Status      string `json:"status"`
		CreatedAt   string `json:"created_at"`
	}

	PaymentMethodResponse struct {
		ID          string `json:"id"`
		PaymentName string `json:"payment_name"`
		CreatedAt   string `json:"created_at"`
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

	GetTransactionResponse struct {
		TransactionId   string            `json:"transactionId,omitempty"`
		TransactionType string            `json:"transactionType,omitempty"`
		UserId          string            `json:"userId,omitempty"`
		UserName        string            `json:"username,omitempty"`
		Amount          string            `json:"amount,omitempty"`
		Description     string            `json:"description"`
		TransactionDate string            `json:"transactionDate"`
		Status          string            `json:"status"`
		TotalDataCount  string            `json:"-"`
		Detail          TransactionDetail `json:"detail"`
	}

	TransactionDetail struct {
		SenderName    string `json:"senderName,omitempty"`
		RecipientName string `json:"recipientName,omitempty"`
		SenderId      string `json:"-"`
		RecipientId   string `json:"-"`
		PaymentMethod string `json:"paymentMethod,omitempty"`
		PaymentURL    string `json:"paymentURL,omitempty"`
		FromWalletId  string `json:"fromWalletId,omitempty"`
		MerchantName  string `json:"merchantName,omitempty"`
		ToWalletId    string `json:"toWalletId,omitempty"`
	}
)
