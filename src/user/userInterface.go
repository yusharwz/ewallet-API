package user

import "final-project-enigma/model/dto/userDto"

type UserRepository interface {
	CekEmail(email string) (bool, error)
	CekPhoneNumber(pnumber string) (bool, error)
	InsertCode(code, email, pnumber string) (bool, error)
	UserLogin(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error)
	UserCreate(req userDto.UserCreateRequest) (userDto.UserCreateResponse, string, error)
	ActivedAccount(req userDto.ActivatedAccountReq) (err error)
	GetDataUserRepo(id string) (userDto.UserGetDataResponse, error)
	GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error)
	GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, error)
	GetTotalDataCount(params userDto.GetTransactionParams) (totalData int, err error)
	UserWalletCreate(id string) (err error)
	CreateTopUpTransaction(req userDto.TopUpTransactionRequest) (string, error)
	CreateWalletTransaction(req userDto.WalletTransactionRequest) (userDto.WalletTransactionResponse, string, error)
	GetPaymentMethodName(id string) (metdhodName string, err error)
	GetUserFullname(id string) (userFullname string, err error)
	PaymentGateway(payload userDto.MidtransSnapReq) (userDto.MidtransSnapResp, error)
	UpdateTransactionStatus(orderID string, status string) error
	UpdateBalance(orderID, amount string) error
	InsertURL(transactionId, url string) error
	UserUploadImage(req userDto.UploadImagesRequest) (userDto.UploadImagesResponse, error)
	ImageToDB(userId string, req userDto.UploadImagesResponse) error
}

type UserUsecase interface {
	LoginCodeReqEmail(email string) error
	LoginCodeReqSMS(pnumber string) error
	LoginReq(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error)
	CreateReq(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error)
	ActivaedAccount(req userDto.ActivatedAccountReq) (err error)
	GetDataUserUC(authHeader string) (userDto.
		UserGetDataResponse, error)
	GetBalanceInfoUC(authHeader string) (resp userDto.UserGetDataResponse, err error)
	GetTransactionUC(authHeader string, params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, string, error)
	TopUpTransaction(req userDto.TopUpTransactionRequest, authHeader string) (userDto.MidtransSnapResp, error)
	WalletTransaction(req userDto.WalletTransactionRequest, authHeader string) (userDto.WalletTransactionResponse, error)
	MidtransStatusReq(notification userDto.MidtransNotification) error
	UploadImagesRequestUC(authHeader string, file userDto.UploadImagesRequest) error
}
