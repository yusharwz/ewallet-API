package router

import (
	"database/sql"

	"final-project-enigma/src/user/userDelivery"
	"final-project-enigma/src/user/userRepository"
	"final-project-enigma/src/user/userUsecase"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB, client *resty.Client) {

	//Users
	userRepo := userRepository.NewUserRepository(db, client)
	userUC := userUsecase.NewUserUsecase(userRepo)
	userDelivery.NewUserDelivery(v1Group, userUC)
}
