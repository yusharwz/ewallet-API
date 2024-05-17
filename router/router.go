package router

import (
	"database/sql"

	"final-project-enigma/src/user/userDelivery"
	"final-project-enigma/src/user/userRepository"
	"final-project-enigma/src/user/userUsecase"

	"github.com/gin-gonic/gin"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB) {

	//Users
	userRepo := userRepository.NewUserRepository(db)
	userUC := userUsecase.NewUserUsecase(userRepo)
	userDelivery.NewUserDelivery(v1Group, userUC)
}
