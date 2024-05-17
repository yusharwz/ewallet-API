package sendEmail

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(email, code string) (bool, error) {
	emailPortStr := os.Getenv("EMAIL_PORT")
	emailPort, _ := strconv.Atoi(emailPortStr)

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_ADDRESS"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/plain", "Here is your verification code: \n"+code)

	d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), emailPort, os.Getenv("EMAIL_ADDRESS"), os.Getenv("EMAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return false, err
	}
	return true, nil
}
