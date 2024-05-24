package auth

import "final-project-enigma/model/dto/userDto"

type AuthRepository interface {
	UserCreate(req userDto.UserCreateRequest) (userDto.UserCreateResponse, string, error)
	UserWalletCreate(id string) (err error)
	ActivedAccount(req userDto.ActivatedAccountReq) error
	CekEmail(email string) (bool, error)
	CekPhoneNumber(pnumber string) (bool, error)
	InsertCode(code, email, pnumber string) (bool, error)
	UserLogin(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error)
	SendLinkForgetPin(req userDto.FogetPinReq) (resp userDto.FogetPinResp, err error)
	ResetPinRepo(req userDto.ForgetPinParams) error
}

type AuthUsecase interface {
	CreateReq(req userDto.UserCreateRequest) (userDto.UserCreateResponse, error)
	ActivatedAccount(req userDto.ActivatedAccountReq) (err error)
	LoginCodeReqEmail(email string) error
	LoginCodeReqSMS(pnumber string) error
	LoginReq(req userDto.UserLoginRequest) (userDto.UserLoginResponse, error)
	ForgotPinReqUC(req userDto.FogetPinReq) error
	ResetPinUC(req userDto.ForgetPinParams) error
}
