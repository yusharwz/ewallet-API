package getJwtToken

import (
	"final-project-enigma/pkg/middleware"
	"fmt"
)

func GetTokenJwt(userId, userEmail string) (string, error) {

	token, err := middleware.GenerateTokenJwt(userId, userEmail, 3)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, nil
}
