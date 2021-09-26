package config

type Search struct {
	ID                   string `yaml:"id"`
	APIKey               string `yaml:"apiKey"`
	CustomSearchEngineId string `yaml:"customSearchEngineId"`
}
