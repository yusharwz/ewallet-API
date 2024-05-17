package user

type UserRepository interface {
	CekEmail(email string) (bool, error)
	InsertCode(code, email string) (bool, error)
}

type UserUsecase interface {
	LoginCodeReq(email string) (string, error)
}
