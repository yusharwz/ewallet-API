package router

import (
	"database/sql"
	"final-project-enigma/src/transaction/transactionRepository"
	"final-project-enigma/src/user/userDelivery"
	"final-project-enigma/src/user/userRepository"
	"final-project-enigma/src/user/userUsecase"

	"final-project-enigma/src/auth/authDelivery"
	"final-project-enigma/src/auth/authRepository"
	"final-project-enigma/src/auth/authUsecase"
	"final-project-enigma/src/payment/paymentDelivery"
	"final-project-enigma/src/payment/paymentRepository"
	"final-project-enigma/src/payment/paymentUsecase"

	"final-project-enigma/src/admin/adminDelivery"
	"final-project-enigma/src/admin/adminRepository"
	"final-project-enigma/src/admin/adminUsecase"

	"final-project-enigma/src/transaction/transactionDelivery"
	"final-project-enigma/src/transaction/transactionUseCase"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB, client *resty.Client) {

	//Auth
	authRepo := authRepository.NewAuthRepository(db)
	authUC := authUsecase.NewAuthUsecase(authRepo)
	authDelivery.NewAuthDelivery(v1Group, authUC)

	//Users
	userRepo := userRepository.NewUserRepository(db, client)
	userUC := userUsecase.NewUserUsecase(userRepo)
	userDelivery.NewUserDelivery(v1Group, userUC)

	//Payment
	paymentRepo := paymentRepository.NewPaymentRepository(db)
	paymentUC := paymentUsecase.NewPaymentUsecase(paymentRepo)
	paymentDelivery.NewPaymentDelivery(v1Group, paymentUC)

	//Admin
	adminRepo := adminRepository.NewAdminRepository(db)
	adminUC := adminUsecase.NewAdminUsecase(adminRepo)
	adminDelivery.NewAdminDelivery(v1Group, adminUC)

	//Transaction
	transactionRepo := transactionRepository.NewTransactionRepository(db)
	transactionUC := transactionUseCase.NewTransactionUseCase(transactionRepo)
	transactionDelivery.NewTransactionDelivery(v1Group, transactionUC)
}
