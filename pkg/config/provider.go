package config

const (
	FREE_MOBILE = "freemobile"
)

type Provider struct {
	ID          string `yaml:"id"`
	Key         string `yaml:"key"`
	Secret      string `yaml:"secret"`
	PhoneNumber string `yaml:"phoneNumber"`
}
