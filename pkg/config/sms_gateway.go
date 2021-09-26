package config

type SMSGateway struct {
	AllowedPhoneNumbers []string `yaml:"allowedPhoneNumbers"`
	Provider            Provider `yaml:"provider"`
}

func (gateway *SMSGateway) IsPhoneNumberAllowed(phoneNumber string) bool {
	for _, allowedPhoneNumber := range gateway.AllowedPhoneNumbers {
		if allowedPhoneNumber == phoneNumber {
			return true
		}
	}
	return false
}
