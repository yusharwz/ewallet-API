package generateCode

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

const (
	expirationDuration = 5 * time.Minute
	timeLayout         = "02/01/2006-15-04-05"
)

func GenerateCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func ValidateCode(encodedTime string) bool {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return false
	}

	decodedTimeStr := string(decodedBytes)
	decodedTime, err := time.Parse(timeLayout, decodedTimeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return false
	}

	if time.Since(decodedTime) > expirationDuration {
		return false
	}

	return true
}
