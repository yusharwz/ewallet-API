package phoneNumberFormat

import "strings"

func ConvertToInternationalFormat(phoneNumber string) string {

	if strings.HasPrefix(phoneNumber, "08") {
		return "+62" + phoneNumber[1:]
	}

	return phoneNumber
}
