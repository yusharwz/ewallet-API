package authUsecase

import (
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/helper/generateCode"
	"final-project-enigma/pkg/helper/getJwtToken"
	"final-project-enigma/pkg/helper/hashingPassword"
	"final-project-enigma/pkg/helper/phoneNumberFormat"
	"final-project-enigma/pkg/helper/sendEmail"
	"final-project-enigma/pkg/helper/sendWhatappTwilio"
	"final-project-enigma/src/auth"

	"github.com/rs/zerolog/log"
)

type authUC struct {
	authRepo auth.AuthRepository
}

func NewAuthUsecase(authRepo auth.AuthRepository) auth.AuthUsecase {
	return &authUC{authRepo}
}

func (usecase *authUC) LoginCodeReqEmail(email string) error {
	resp, err := usecase.authRepo.CekEmail(email)
	if err != nil {
		return err
	}
	if resp.Status != "active" {
		code, err := generateCode.GenerateCode()
		if err != nil {
			return err
		}

		var pnumber string
		respInsertCode, err := usecase.authRepo.InsertCode(code, resp.Email, pnumber)
		if err != nil {
			return err
		}
		if !respInsertCode {
			return err
		}

		err = sendEmail.SendEmailActivedAccount(resp.Email, resp.Username, code, resp.Unique)
		if err != nil {
			return err
		}
		log.Error().Msg("account has not been activated, please check the email inbox for the activation link")
		return errors.New("account has not been activated, please check the email inbox for the activation link")
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}

	var pnumber string
	respInsertCode, err := usecase.authRepo.InsertCode(code, email, pnumber)
	if err != nil {
		return err
	}
	if !respInsertCode {
		return err
	}

	emailResp, err := sendEmail.SendEmail(email, code)

	if err != nil {
		log.Error().Msg("failed to send email")
		return errors.New("failed to send email")
	}

	if !emailResp {
		log.Error().Msg("failed to send email")
		return errors.New("failed to send email")
	}
	return nil
}

func (usecase *authUC) LoginCodeReqSMS(pnumber string) error {

	phoneNumberFormatted := phoneNumberFormat.ConvertToInternationalFormat(pnumber)
	pnumber = phoneNumberFormatted

	resp, err := usecase.authRepo.CekPhoneNumber(pnumber)
	if err != nil {
		return err
	}

	if resp.Status != "active" {
		code, err := generateCode.GenerateCode()
		if err != nil {
			return err
		}

		var pnumber string
		respInsertCode, err := usecase.authRepo.InsertCode(code, resp.Email, pnumber)
		if err != nil {
			return err
		}
		if !respInsertCode {
			return err
		}

		err = sendEmail.SendEmailActivedAccount(resp.Email, resp.Username, code, resp.Unique)
		if err != nil {
			return err
		}
		log.Error().Msg("account has not been activated, please check the email inbox for the activation link")
		return errors.New("account has not been activated, please check the email inbox for the activation link")
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}

	var email string
	respInsertCode, err := usecase.authRepo.InsertCode(code, email, pnumber)
	if err != nil {
		return err
	}
	if !respInsertCode {
		return err
	}

	emailResp, err := sendWhatappTwilio.SendWhatsAppMessage(pnumber, code)

	if err != nil {
		return err
	}

	if !emailResp {
		return err
	}

	return nil
}

func (usecase *authUC) LoginReq(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {

	resp, err = usecase.authRepo.UserLogin(req)
	if err != nil {
		return resp, err
	}

	err = hashingPassword.ComparePassword(resp.Pin, req.Pin)
	if err != nil {
		log.Error().Msg("invalid pin or verification code")
		return resp, errors.New("invalid pin or verification code")
	}

	resp.Token, err = getJwtToken.GetTokenJwt(resp.UserId, resp.UserEmail, resp.Roles)
	if err != nil {
		return resp, err
	}

	resp.UserId = ""
	resp.UserEmail = ""
	resp.Pin = ""

	return resp, nil
}

func (usecase *authUC) CreateReq(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, err error) {

	hashedPin, err := hashingPassword.HashPassword(req.Pin)
	if err != nil {
		return resp, err
	}

	phoneNumberFormatted := phoneNumberFormat.ConvertToInternationalFormat(req.PhoneNumber)
	req.PhoneNumber = phoneNumberFormatted

	req.Pin = hashedPin
	if req.Roles == "" {
		req.Roles = "USER"
	}

	resp, unique, err := usecase.authRepo.UserCreate(req)
	if err != nil {
		return resp, err
	}

	err = usecase.authRepo.UserWalletCreate(resp.Id)
	if err != nil {
		return resp, err
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return resp, err
	}

	var pnumber string
	respInsertCode, err := usecase.authRepo.InsertCode(code, resp.Email, pnumber)
	if err != nil {
		return resp, err
	}
	if !respInsertCode {
		return resp, err
	}

	err = sendEmail.SendEmailActivedAccount(resp.Email, resp.Username, code, unique)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (usecase *authUC) ActivatedAccount(req userDto.ActivatedAccountReq) (err error) {

	err = usecase.authRepo.ActivedAccount(req)
	if err != nil {
		return err
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}
	var pnumber string
	respInsertCode, err := usecase.authRepo.InsertCode(code, req.Email, pnumber)
	if err != nil {
		return err
	}
	if !respInsertCode {
		return err
	}

	return nil
}

func (usecase *authUC) ForgotPinReqUC(req userDto.ForgetPinReq) (err error) {

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}
	var pnumber string
	respInsertCode, err := usecase.authRepo.InsertCode(code, req.Email, pnumber)
	if err != nil {
		return err
	}
	if !respInsertCode {
		return err
	}

	resp, err := usecase.authRepo.SendLinkForgetPin(req)
	if err != nil {
		return err
	}

	err = sendEmail.SendEmailForgotPin(resp.Email, resp.Username, resp.Code, resp.Unique)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *authUC) ResetPinUC(req userDto.ForgetPinParams) error {

	if req.NewPin != req.RetypeNewPin {
		log.Error().Msg("new pin and retype new pin not match")
		return errors.New("new pin and retype new pin not match")
	}

	hashedPin, err := hashingPassword.HashPassword(req.NewPin)
	if err != nil {
		return err
	}

	req.NewPin = hashedPin

	err = usecase.authRepo.ResetPinRepo(req)
	if err != nil {
		return err
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}
	var pnumber string
	respInsertCode, err := usecase.authRepo.InsertCode(code, req.Email, pnumber)
	if err != nil {
		return err
	}
	if !respInsertCode {
		return err
	}

	return nil
}
