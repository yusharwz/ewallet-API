package userDelivery

import (
	"final-project-enigma/model/dto/json"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/pkg/validation"
	"final-project-enigma/src/user"

	"github.com/gin-gonic/gin"
)

type authorDelivery struct {
	userUC user.UserUsecase
}

func NewUserDelivery(v1Group *gin.RouterGroup, userUC user.UserUsecase) {
	handler := authorDelivery{
		userUC: userUC,
	}
	userGroup := v1Group.Group("/users")

	{
		userGroup.POST("/login", middleware.BasicAuth, handler.loginUserCodeReuqest)
	}
}

func (u *authorDelivery) loginUserCodeReuqest(ctx *gin.Context) {
	var req userDto.UserLoginCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)

		if len(validationError) > 0 {
			json.NewResponBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
		json.NewResponseError(ctx, "json request body required", "01", "02")
		return
	}

	resp, err := u.userUC.LoginCodeReq(req.Email)
	if err != nil {
		json.NewResponseForbidden(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponSucces(ctx, "email: "+req.Email, resp, "01", "01")
}
