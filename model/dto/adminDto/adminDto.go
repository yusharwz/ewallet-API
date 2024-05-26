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
	GetpaymentMethodParams struct {
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
)
