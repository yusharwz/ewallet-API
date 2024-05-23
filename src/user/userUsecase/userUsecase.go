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
	"strconv"
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
		return err
	}

	if !result {
		return err
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}

	var pnumber string
	respInsertCode, err := usecase.userRepo.InsertCode(code, email, pnumber)
	if err != nil {
		return err
	}
	if !respInsertCode {
		return err
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
		return err
	}

	if !result {
		return err
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return err
	}

	var email string
	respInsertCode, err := usecase.userRepo.InsertCode(code, email, pnumber)
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

func (usecase *userUC) LoginReq(req userDto.UserLoginRequest) (resp userDto.UserLoginResponse, err error) {

	resp, err = usecase.userRepo.UserLogin(req)
	if err != nil {
		return resp, err
	}

	err = hashingPassword.ComparePassword(resp.Pin, req.Pin)
	if err != nil {
		return resp, err
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

func (usecase *userUC) UploadImagesRequestUC(authHeader string, file userDto.UploadImagesRequest) error {

	userId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return err
	}

	resp, err := usecase.userRepo.UserUploadImage(file)
	if err != nil {
		return err
	}

	err = usecase.userRepo.ImageToDB(userId, resp)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *userUC) CreateReq(req userDto.UserCreateRequest) (resp userDto.UserCreateResponse, err error) {

	hashedPin, err := hashingPassword.HashPassword(req.Pin)
	if err != nil {
		return resp, err
	}

	req.Pin = hashedPin
	req.Roles = "USER"

	resp, unique, err := usecase.userRepo.UserCreate(req)
	if err != nil {
		return resp, err
	}

	err = usecase.userRepo.UserWalletCreate(resp.Id)
	if err != nil {
		return resp, err
	}

	code, err := generateCode.GenerateCode()
	if err != nil {
		return resp, err
	}

	var pnumber string
	respInsertCode, err := usecase.userRepo.InsertCode(code, resp.Email, pnumber)
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

func (usecase *userUC) ActivaedAccount(req userDto.ActivatedAccountReq) (err error) {

	err = usecase.userRepo.ActivedAccount(req)
	if err != nil {
		return err
	}

	return nil
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
	return resp, nil
}

func (usecase *userUC) GetTransactionUC(authHeader string, params userDto.GetTransactionParams) ([]userDto.GetTransactionResponse, string, error) {
	userId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return nil, "", err
	}

	params.UserId = userId

	resp, err := usecase.userRepo.GetTransactionRepo(params)
	if err != nil {
		return nil, "", err
	}
	for i := range resp {
		if resp[i].Detail.RecipientId == params.UserId {
			resp[i].TransactionType = "credit"
		}
	}

	totalData, err := usecase.userRepo.GetTotalDataCount(params)
	if err != nil {
		return nil, "", err
	}

	totalDataStr := strconv.Itoa(totalData)
	return resp, totalDataStr, nil
}

func (usecase *userUC) TopUpTransaction(req userDto.TopUpTransactionRequest, authHeader string) (userDto.MidtransSnapResp, error) {
	userId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return userDto.MidtransSnapResp{}, err
	}

	req.UserId = userId
	req.Description = "Balance Top Up"

	transactionId, err := usecase.userRepo.CreateTopUpTransaction(req)
	if err != nil {
		return userDto.MidtransSnapResp{}, err
	}

	methodName, err := usecase.userRepo.GetPaymentMethodName(req.PaymentMethodId)
	if err != nil {
		return userDto.MidtransSnapResp{}, err
	}
	userFullname, err := usecase.userRepo.GetUserFullname(req.UserId)
	if err != nil {
		return userDto.MidtransSnapResp{}, err
	}
	var request userDto.MidtransSnapReq
	request.TransactionDetail.OrderID = transactionId
	request.TransactionDetail.GrossAmt = req.Amount

	request.PaymentType = methodName

	request.Customer = userFullname

	items := []userDto.Item{
		{
			ID:       "usertopup",
			Name:     "TopUp Balance",
			Price:    req.Amount,
			Quantity: 1,
		},
	}

	request.Items = items

	resp, err := usecase.userRepo.PaymentGateway(request)

	if err := usecase.userRepo.InsertURL(transactionId, resp.RedirectUrl); err != nil {
		return userDto.MidtransSnapResp{}, err
	}
	qrisPath := "#/other-qris"
	bcaVaPath := "#/bank-transfer/bca-va"
	mandiriVaPath := "#/bank-transfer/mandiri-va"
	bniVaPath := "#/bank-transfer/bni-va"
	briVaPath := "#/bank-transfer/bri-va"
	permataVaPath := "#/bank-transfer/permata-va"
	cimbnVaPath := "#/bank-transfer/cimb-va"
	gopayPath := "#/gopay-qris"
	debitCreditCardPath := "#/credit-card"
	spaySpayLaterPath := "#/shopeepay-qris"
	alfamartPath := "#/alfamart"
	indomaretPath := "#/indomaret"
	akulakuPath := "#/akulaku"
	kredivoPath := "#/kredivo"

	if req.PaymentMethodId == "089e8004-2428-41f9-bf06-856082bb83d3" {
		resp.RedirectUrl += qrisPath
	}
	if req.PaymentMethodId == "cf51fa64-1686-4fee-a4e1-ea13c939f99b" {
		resp.RedirectUrl += bcaVaPath
	}
	if req.PaymentMethodId == "087f9751-1dfc-474d-bdee-07ce44b1fe7a" {
		resp.RedirectUrl += mandiriVaPath
	}
	if req.PaymentMethodId == "2bed0329-499e-43b5-9b99-583b203ea102" {
		resp.RedirectUrl += bniVaPath
	}
	if req.PaymentMethodId == "3863b99e-9909-486c-8ec1-b7a3162c9f97" {
		resp.RedirectUrl += briVaPath
	}
	if req.PaymentMethodId == "76954351-6cb3-496d-8866-d7f5772a04fe" {
		resp.RedirectUrl += permataVaPath
	}
	if req.PaymentMethodId == "0fafc78f-ebbf-421d-bc89-3246ce6198ad" {
		resp.RedirectUrl += cimbnVaPath
	}
	if req.PaymentMethodId == "f9569b06-a389-4685-b3cc-89b13a111214" {
		resp.RedirectUrl += gopayPath
	}
	if req.PaymentMethodId == "9fa520e0-d10b-4be1-a6d7-e8b6fc635c5c" {
		resp.RedirectUrl += debitCreditCardPath
	}
	if req.PaymentMethodId == "91b75dee-155e-4ac3-9bfd-f8bed82b6189" {
		resp.RedirectUrl += spaySpayLaterPath
	}
	if req.PaymentMethodId == "b25a226e-82ab-4d29-a68e-6957fb7e21a9" {
		resp.RedirectUrl += alfamartPath
	}
	if req.PaymentMethodId == "0eaad501-e44d-46e2-902a-9325c6c6c5eb" {
		resp.RedirectUrl += indomaretPath
	}
	if req.PaymentMethodId == "29690f9f-c6c4-4fda-acac-be91555b1f94" {
		resp.RedirectUrl += akulakuPath
	}
	if req.PaymentMethodId == "220309af-cd3b-40e5-b353-6754c66f3831" {
		resp.RedirectUrl += kredivoPath
	}

	return resp, err
}

func (usecase *userUC) WalletTransaction(req userDto.WalletTransactionRequest, authHeader string) (userDto.WalletTransactionResponse, error) {
	fromId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return userDto.WalletTransactionResponse{}, err
	}
	req.UserId = fromId

	resp, storedPin, err := usecase.userRepo.CreateWalletTransaction(req)
	if err != nil {
		return userDto.WalletTransactionResponse{}, err
	}

	err = hashingPassword.ComparePassword(storedPin, req.PIN)
	if err != nil {
		return userDto.WalletTransactionResponse{}, errors.New("invalid PIN")
	}
	return resp, nil
}

func (usecase *userUC) MidtransStatusReq(notification userDto.MidtransNotification) error {

	switch notification.TransactionStatus {
	case "capture":
		if notification.FraudStatus == "challenge" {
		} else if notification.FraudStatus == "accept" {
			usecase.userRepo.UpdateTransactionStatus(notification.OrderID, "settlement")
		}
	case "settlement":
		if err := usecase.userRepo.UpdateTransactionStatus(notification.OrderID, "succes"); err != nil {
			return err
		}
		if err := usecase.userRepo.UpdateBalance(notification.OrderID, notification.GrossAmount); err != nil {
			return err
		}
	case "deny":
		if err := usecase.userRepo.UpdateTransactionStatus(notification.OrderID, "deny"); err != nil {
			return err
		}
	case "cancel", "expire":
		if err := usecase.userRepo.UpdateTransactionStatus(notification.OrderID, "cancel"); err != nil {
			return err
		}
	case "pending":
		if err := usecase.userRepo.UpdateTransactionStatus(notification.OrderID, "pending"); err != nil {
			return err
		}
	default:
		usecase.userRepo.UpdateTransactionStatus(notification.OrderID, notification.TransactionStatus)
	}

	return nil
}
