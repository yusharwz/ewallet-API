package userDto

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

	GetTransactionResponse struct {
		TransactionId   string            `json:"transactionId,omitempty"`
		TransactionType string            `json:"transactionType,omitempty"`
		Amount          string            `json:"amount,omitempty"`
		Description     string            `json:"description"`
		TransactionDate string            `json:"transactionDate"`
		Status          string            `json:"status"`
		Detail          TransactionDetail `json:"detail"`
	}

	TransactionDetail struct {
		SenderName    string `json:"senderName,omitempty"`
		RecipientName string `json:"recipientName,omitempty"`
		PaymentMethod string `json:"paymentMethod,omitempty"`
	}

	TopUpTransactionRequest struct {
		UserId          string  `json:"userId"`
		Amount          float64 `json:"amount"`
		Description     string  `json:"description"`
		PaymentMethodId string  `json:"paymentMethodId"`
	}

	TopUpTransactionResponse struct {
		TransactionId string `json:"transactionId"`
	}

	WalletTransactionRequest struct {
		UserId       string  `json:"userId"`
		FromWalletId string  `json:"fromWalletId"`
		ToWalletId   string  `json:"toWalletId"`
		Amount       float64 `json:"amount"`
		Description  string  `json:"description"`
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
)
