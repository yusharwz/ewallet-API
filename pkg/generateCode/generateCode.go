package generateCode

import (
	"encoding/base64"
	"fmt"
	"time"
)

const (
	expirationDuration = 5 * time.Minute
	timeLayout         = "02/01/2006-15-04-05"
)

func GenerateCode() string {
	// Ambil timestamp saat ini
	currentTime := time.Now()

	// Format timestamp ke dalam format yang diinginkan: dd/mm/yyyy-hh-mm-ss
	formattedTime := currentTime.Format(timeLayout)

	// Encode timestamp yang diformat ke dalam base64
	encodedTime := base64.StdEncoding.EncodeToString([]byte(formattedTime))

	return encodedTime
}

func ValidateCode(encodedTime string) bool {
	// Decode base64 menjadi timestamp yang diformat
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return false
	}

	// Parse timestamp dari string yang di-decode
	decodedTimeStr := string(decodedBytes)
	decodedTime, err := time.Parse(timeLayout, decodedTimeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return false
	}

	// Cek apakah timestamp masih berlaku (kurang dari 5 menit)
	if time.Since(decodedTime) > expirationDuration {
		return false
	}

	return true
}
