package getJwtToken

import (
	"errors"
	"final-project-enigma/pkg/middleware"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetTokenJwt(userId, userEmail, reqPassword, hashedPassword string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(reqPassword))
	if err != nil {
		return "", errors.New("Password salah")
	}
	token, err := middleware.GenerateTokenJwt(userId, userEmail, 3)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, nil
}
