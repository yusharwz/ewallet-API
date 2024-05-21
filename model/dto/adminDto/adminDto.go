package adminDto

import "time"

type (
	GetUserParams struct {
		ID          string `json:"id"`
		Fullname    string `json:"fullname"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		CreateAt    string `json:"createAt"`
		Page        string `json:"page"`
		Limit       string `json:"limit"`
	}
	User struct {
		ID          string    `json:"id"`
		Fullname    string    `json:"fullname"`
		Username    string    `json:"username"`
		Pin         string    `json:"pin"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phoneNumber"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		DeletedAt   time.Time `json:"deletedAt"`
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
		DeletedAt   time.Time `json:"deletedAt"`
	}

	GetWalletParams struct {
		ID         string   `json:"id"`
		User_id    string   `json:"user_id"`
		MinBalance *float64 `json:"min_balance"`
		MaxBalance *float64 `json:"max_balance"`
		CreatedAt  string   `json:"createdAt"`
		Page       string   `json:"page"`
		Limit      string   `json:"limit"`
	}
	Wallet struct {
		ID        string    `json:"id"`
		User_id   string    `json:"user_id"`
		Balance   string    `json:"balance"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)
