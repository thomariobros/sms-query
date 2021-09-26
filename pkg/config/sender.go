package config

type Sender struct {
	PhoneNumber string   `yaml:"phoneNumber"`
	Provider    Provider `yaml:"provider"`
}
