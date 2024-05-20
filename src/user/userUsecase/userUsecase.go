package userUsecase

import (
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/generateCode"
	"final-project-enigma/pkg/getJwtToken"
	"final-project-enigma/pkg/hashingPassword"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/pkg/sendEmail"
	"final-project-enigma/pkg/sendWhatappTwilio"
	"final-project-enigma/src/user"
	"fmt"
)

type userUC struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUC{userRepo}
}

func (usecase *userUC) LoginCodeReqEmail(email string) error {
	result, err := usecase.userRepo.CekEmail(email)
	if err != nil {
		return errors.New("Email not found: " + email)
	}

	if !result {
		return errors.New("Email not found: " + email)
	}

	code := generateCode.GenerateCode()

	var pnumber string
	respInsertCode, err := usecase.userRepo.InsertCode(code, email, pnumber)
	if err != nil {
		return errors.New("failed to insert code")
	}
	if !respInsertCode {
		return errors.New("failed to insert code")
	}

	emailResp, err := sendEmail.SendEmail(email, code)

	if err != nil {
		return errors.New("failed to send email")
	}

	if !emailResp {
		return errors.New("failed to send email")
	}
	return nil
}

func (usecase *userUC) LoginCodeReqSMS(pnumber string) error {
	result, err := usecase.userRepo.CekPhoneNumber(pnumber)
	if err != nil {
		return errors.New("phone number not found")
	}

	if !result {
		return errors.New("phone number not found")
	}

	code := generateCode.GenerateCode()

	var email string
	respInsertCode, err := usecase.userRepo.InsertCode(code, email, pnumber)
	if err != nil {
		return errors.New("fail to insert code")
	}
	if !respInsertCode {
		return errors.New("fail to insert code")
	}

	emailResp, err := sendWhatappTwilio.SendWhatsAppMessage(pnumber, code)

	if err != nil {
		return errors.New("fail to send email")
	}

	if !emailResp {
		return errors.New("fail to send email")
	}

	return nil
}

func (usecase *userUC) LoginReq(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {

	resp, err = usecase.userRepo.UserLogin(req)
	if err != nil {
		return resp, err
	}

	err = hashingPassword.ComparePassword(resp.Pin, req.Pin)
	if err != nil {
		return resp, err
	}

	resp.Token, err = getJwtToken.GetTokenJwt(resp.UserId, resp.UserEmail)
	if err != nil {
		return resp, err
	}

	resp.UserId = ""
	resp.UserEmail = ""
	resp.Pin = ""

	return resp, nil
}

func (usecase *userUC) CreateReq(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, err error) {

	hashedPin, err := hashingPassword.HashPassword(req.Pin)
	if err != nil {
		return resp, err
	}

	req.Pin = hashedPin

	resp, err = usecase.userRepo.UserCreate(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (usecase *userUC) GetDataUserUC(authHeader string) (resp userDto.UserGetDataResponse, err error) {

	id, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return resp, err
	}

	resp, err = usecase.userRepo.GetDataUserRepo(id)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (usecase *userUC) GetBalanceInfoUC(authHeader string) (resp userDto.UserGetDataResponse, err error) {

	id, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return resp, err
	}

	resp, err = usecase.userRepo.GetBalanceInfoRepo(id)
	if err != nil {
		return resp, err
	}
	fmt.Println(resp.Balance)
	return resp, nil
}

func (usecase *userUC) GetTransactionUC(authHeader string, params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, error) {
	userId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return nil, err
	}

	params.UserId = userId

	transactions, err := usecase.userRepo.GetTransactionRepo(params)
	if err != nil {
		return nil, err
	}

	var resp []userDto.GetTransactionResponse
	for _, transaction := range transactions {
		trxType := "debit"
		if transaction.PaymentMethod.Valid && transaction.PaymentMethod.String != "" {
			trxType = "credit"
		} else if transaction.RecipientUserId.Valid && transaction.RecipientUserId.String == userId {
			trxType = "credit"
		}

		if params.TrxType != "" && trxType != params.TrxType {
			continue
		}

		resp = append(resp, userDto.GetTransactionResponse{
			TransactionId:   transaction.TransactionId,
			PaymentMethod:   transaction.PaymentMethod.String,
			TransactionType: trxType,
			Amount:          transaction.Amount,
			Description:     transaction.Description,
			TransactionDate: transaction.TransactionDate,
			PaymentStatus:   transaction.PaymentStatus,
			SenderName:      transaction.SenderName.String,
			RecipientName:   transaction.RecipientName.String,
		})
	}

	return resp, nil
}
