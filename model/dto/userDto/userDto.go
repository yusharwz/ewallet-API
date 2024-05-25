package userDto

import (
	"mime/multipart"
)

type (
	UserLoginCodeRequestEmail struct {
		Email string `json:"email" binding:"required,email"`
	}

	UserLoginCodeRequestPhoneNumber struct {
		PhoneNumber string `json:"phoneNumber" binding:"required,nomorHp"`
	}

	UserLoginRequest struct {
		Email string `json:"email" binding:"required,email"`
		Pin   string `json:"pin" binding:"required,pin,min=6,max=6"`
		Code  string `json:"code" binding:"required"`
	}

	UploadImagesRequest struct {
		File multipart.File
	}

	UploadImagesResponse struct {
		Url string
	}

	UserLoginResponse struct {
		UserId    string `json:"userId,omitempty"`
		UserEmail string `json:"userEmail,omitempty"`
		Pin       string `json:"pin,omitempty"`
		Token     string `json:"token,omitempty"`
		Roles     string `json:"roles,omitempty"`
		Status    string `json:"status,omitempty"`
	}

	UserUpdateReq struct {
		UserId      string
		Fullname    string `json:"fullname"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}

	UserCreateRequest struct {
		Fullname    string `json:"fullname" binding:"required,min=1"`
		Username    string `json:"username" binding:"required,username,min=5,max=20"`
		Email       string `json:"email" binding:"required,email"`
		Pin         string `json:"pin" binding:"required,pin,min=6,max=6"`
		PhoneNumber string `json:"phoneNumber" binding:"required,nomorHp,min=8,max=17"`
		Roles       string
	}

	ActivatedAccountReq struct {
		Email    string
		Fullname string
		Unique   string
		Code     string
	}

	UserCreateResponse struct {
		Id          string `json:"id"`
		Fullname    string `json:"fullname"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}

	ForgetPinReq struct {
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber string `json:"phoneNumber" binding:"required,nomorHp"`
	}

	ForgetPinResp struct {
		Email    string
		Username string
		Code     string
		Unique   string
		Status   string
	}

	ForgetPinParams struct {
		Email        string
		NewPin       string `json:"newPin" binding:"required,pin"`
		RetypeNewPin string `json:"retypeNewPin" binding:"required,pin"`
		Username     string
		Code         string
		Unique       string
	}

	UserGetDataResponse struct {
		Fullname     string `json:"fullname,omitempty"`
		Username     string `json:"username,omitempty"`
		Email        string `json:"email,omitempty"`
		ProfilImages string `json:"profilImages,omitempty"`
		PhoneNumber  string `json:"phoneNumber,omitempty"`
		Balance      string `json:"balance,omitempty"`
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
		ToWalletId    string `json:"toWalletId,omitempty"`
	}

	TopUpTransactionRequest struct {
		UserId          string  `json:"userId"`
		Amount          float64 `json:"amount" binding:"required,min=5"`
		Description     string  `json:"description"`
		PaymentMethodId string  `json:"paymentMethodId" binding:"required,min=15"`
	}

	TopUpTransactionResponse struct {
		TransactionId string `json:"transactionId"`
	}

	WalletTransactionRequest struct {
		UserId               string `json:"userId"`
		FromWalletId         string `json:"fromWalletId"`
		ToWalletId           string `json:"toWalletId"`
		RecipientPhoneNumber string `json:"recipientPhoneNumber" binding:"required"`
		Amount               string `json:"amount" binding:"required,min=5"`
		PIN                  string `json:"pin" binding:"required,pin"`
		Description          string `json:"description"`
	}

	WalletTransactionResponse struct {
		TransactionId string `json:"transactionId"`
	}

	MidtransSnapReq struct {
		TransactionDetail struct {
			OrderID  string  `json:"order_id"`
			GrossAmt float64 `json:"gross_amount"`
		} `json:"transaction_details"`
		PaymentType string `json:"payment_type"`
		Customer    string `json:"customer"`
		Items       []Item `json:"item_details"`
	}

	Item struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	MidtransSnapResp struct {
		Token        string   `json:"token"`
		RedirectUrl  string   `json:"redirect_url"`
		ErrorMessage []string `json:"error_messages"`
	}

	MidtransNotification struct {
		TransactionStatus string `json:"transaction_status"`
		OrderID           string `json:"order_id"`
		GrossAmount       string `json:"gross_amount"`
		PaymentType       string `json:"payment_type"`
		TransactionTime   string `json:"transaction_time"`
		FraudStatus       string `json:"fraud_status"`
	}
)
