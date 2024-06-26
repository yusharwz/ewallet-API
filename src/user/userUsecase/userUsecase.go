package userUsecase

import (
	"errors"
	"final-project-enigma/model/dto/userDto"
	"final-project-enigma/pkg/helper/hashingPassword"
	"final-project-enigma/pkg/middleware"
	"final-project-enigma/src/user"
	"strconv"

	"github.com/rs/zerolog/log"
)

type userUC struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUC{userRepo}
}

func (usecase *userUC) EditDataUserUC(authHeader string, req userDto.UserUpdateReq) error {

	userId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return err
	}

	req.UserId = userId

	currentUserData, err := usecase.userRepo.GetDataUserRepo(req.UserId)
	if err != nil {
		return err
	}

	if req.Fullname == "" {
		req.Fullname = currentUserData.Fullname
	}
	if req.Email == "" {
		req.Email = currentUserData.Email
	}
	if req.PhoneNumber == "" {
		req.PhoneNumber = currentUserData.PhoneNumber
	}

	if err := usecase.userRepo.EditUserData(req); err != nil {
		return err
	}

	return nil
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

	resp, totalData, err := usecase.userRepo.GetTransactionRepo(params)
	if err != nil {
		return nil, "", err
	}

	var manipulatedTransactions []userDto.GetTransactionResponse
	for _, transaction := range resp {
		if transaction.Detail.RecipientId == userId {
			transaction.TransactionType = "credit"
		} else if transaction.Detail.PaymentMethod != "" {
			transaction.TransactionType = "credit"
		} else if transaction.Detail.MerchantName != "" {
			transaction.TransactionType = "debit"
		} else {
			transaction.TransactionType = "debit"
		}

		manipulatedTransactions = append(manipulatedTransactions, transaction)
	}

	if params.TrxType != "" {
		var filteredTransactions []userDto.GetTransactionResponse
		for _, trx := range manipulatedTransactions {
			if trx.TransactionType == params.TrxType {
				filteredTransactions = append(filteredTransactions, trx)
			}
		}
		manipulatedTransactions = filteredTransactions
		totalData = len(filteredTransactions)
	}

	totalDataStr := strconv.Itoa(totalData)
	return manipulatedTransactions, totalDataStr, nil
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

	if err := usecase.userRepo.InsertPaymentURL(transactionId, resp.RedirectUrl); err != nil {
		return userDto.MidtransSnapResp{}, err
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
		log.Error().Msg("invalid PIN")
		return userDto.WalletTransactionResponse{}, errors.New("invalid PIN")
	}
	return resp, nil
}

func (usecase *userUC) DeleteUser(authHeader string) error {
	id, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return err
	}
	return usecase.userRepo.DeleteUser(id)
}

func (usecase *userUC) MerchantTransaction(req userDto.MerchantTransactionRequest, authHeader string) (resp userDto.MerchantTransactionResponse, err error) {
	userId, err := middleware.GetIdFromToken(authHeader)
	if err != nil {
		return userDto.MerchantTransactionResponse{}, err
	}

	req.UserId = userId
	req.Description = "Merchant-Payment"

	transactionId, err := usecase.userRepo.CreateMerchantTransaction(req)
	if err != nil {
		return userDto.MerchantTransactionResponse{}, err
	}

	resp.TransactionId = transactionId

	return resp, nil
}
