package getJwtToken

import (
	"final-project-enigma/pkg/middleware"
	"fmt"
)

func GetTokenJwt(userId, userEmail, roles string) (string, error) {

	token, err := middleware.GenerateTokenJwt(userId, userEmail, roles, 720)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, nil
}
