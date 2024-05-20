package sendWhatappTwilio

import (
	"fmt"
	"os"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendWhatsAppMessage(to, code string) (bool, error) {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_WHATSAPP_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	expiredCode := time.Now().Add(5 * time.Minute)

	message := fmt.Sprintf("Here is your verification code: \n%s \n your code expired on 5 minutes (%s)", code, expiredCode)
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return false, err
	}

	return true, nil
}
