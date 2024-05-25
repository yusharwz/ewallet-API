package user

import "final-project-enigma/model/dto/userDto"

type UserRepository interface {
	UserUploadImage(req userDto.UploadImagesRequest) (userDto.UploadImagesResponse, error)
	ImageToDB(userId string, req userDto.UploadImagesResponse) error
	GetDataUserRepo(id string) (userDto.UserGetDataResponse, error)
	GetBalanceInfoRepo(id string) (resp userDto.UserGetDataResponse, err error)
	GetTransactionRepo(params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, int, error)
	CreateTopUpTransaction(req userDto.TopUpTransactionRequest) (string, error)
	GetPaymentMethodName(id string) (metdhodName string, err error)
	GetUserFullname(id string) (userFullname string, err error)
	PaymentGateway(payload userDto.MidtransSnapReq) (userDto.MidtransSnapResp, error)
	InsertPaymentURL(transactionId, url string) error
	CreateWalletTransaction(req userDto.WalletTransactionRequest) (userDto.WalletTransactionResponse, string, error)
	EditUserData(req userDto.UserUpdateReq) error
	DeleteUser(id string) error
}

type UserUsecase interface {
	UploadImagesRequestUC(authHeader string, file userDto.UploadImagesRequest) error
	GetDataUserUC(authHeader string) (userDto.
		UserGetDataResponse, error)
	GetBalanceInfoUC(authHeader string) (resp userDto.UserGetDataResponse, err error)
	GetTransactionUC(authHeader string, params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, string, error)
	TopUpTransaction(req userDto.TopUpTransactionRequest, authHeader string) (userDto.MidtransSnapResp, error)
	WalletTransaction(req userDto.WalletTransactionRequest, authHeader string) (userDto.WalletTransactionResponse, error)
	EditDataUserUC(authHeader string, req userDto.UserUpdateReq) error
	DeleteUser(authHeader string) error
}
