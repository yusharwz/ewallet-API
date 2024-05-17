package authorDto

type (
	AuthorLoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	AuthorLoginResponse struct {
		Token string `json:"token"`
	}
)
