package userDto

type (
	UserLoginCodeRequest struct {
		Email string `json:"email" binding:"required,email"`
	}
)
