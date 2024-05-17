package userUsecase

import (
	"final-project-enigma/pkg/generateCode"
	"final-project-enigma/pkg/sendEmail"
	"final-project-enigma/src/user"
)

type userUC struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUC{userRepo}
}

func (usecase *userUC) LoginCodeReq(email string) (string, error) {
	resp, err := usecase.userRepo.CekEmail(email)
	if err != nil {
		return "email tidak terdaftar", err
	}

	if !resp {
		return "email tidak terdaftar", err
	}

	code := generateCode.GenerateCode()

	respInsertCode, err := usecase.userRepo.InsertCode(code, email)
	if err != nil {
		return "Fail to insert code", err
	}
	if !respInsertCode {
		return "Fail to insert code", err
	}

	emailResp, err := sendEmail.SendEmail(email, code)

	if err != nil {
		return "Fail to send email", err
	}

	if !emailResp {
		return "Fail to send email", err
	}

	return "Please check your email inbox", nil
}
