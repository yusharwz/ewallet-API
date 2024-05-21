package user

import "final-project-enigma/model/dto/userDto"

type UserRepository interface {
	CekEmail(email string) (bool, error)
	CekPhoneNumber(pnumber string) (bool, error)
	InsertCode(code, email, pnumber string) (bool, error)
	UserLogin(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error)
	UserCreate(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error)
	GetDataUserRepo(id string) (userDto.UserGetDataResponse, error)
	GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error)
	GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, error)
	UserWalletCreate(id string) (err error)
	CreateTopUpTransaction(req userDto.TopUpTransactionRequest) (userDto.TopUpTransactionResponse, error)
	CreateWalletTransaction(req userDto.WalletTransactionRequest) (userDto.WalletTransactionResponse, error)
}

type UserUsecase interface {
	LoginCodeReqEmail(email string) error
	LoginCodeReqSMS(pnumber string) error
	LoginReq(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error)
	CreateReq(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error)
	GetDataUserUC(authHeader string) (userDto.
		UserGetDataResponse, error)
	GetBalanceInfoUC(authHeader string) (resp userDto.UserGetDataResponse, err error)
	GetTransactionUC(authHeader string, params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, error)
	TopUpTransaction(req userDto.TopUpTransactionRequest, authHeader string) (userDto.TopUpTransactionResponse, error)
	WalletTransaction(req userDto.WalletTransactionRequest, authHeader string) (userDto.WalletTransactionResponse, error)
}
