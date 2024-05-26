package hashingPassword

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

type Hasher interface {
	HashPassword(password string) (string, error)
}


func HashPasswordMock(password string) (string, error) {
	// Implementasi hashing password
	return "hashedPassword", nil
}

type BcryptHasher struct{}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
    return HashPassword(password) // fungsi asli
}