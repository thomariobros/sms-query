package config

const (
	FREE_MOBILE = "freemobile"
)

type Provider struct {
	ID              string `yaml:"id"`
	Key             string `yaml:"key"`
	Secret          string `yaml:"secret"`
	SignatureSecret string `yaml:"signatureSecret"`
	PhoneNumber     string `yaml:"phoneNumber"`
}
