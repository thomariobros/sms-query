package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

var instance *Config
var once sync.Once

func GetInstance() *Config {
	return instance
}

func InitWithRootPath(rootPath string) error {
	var err error
	once.Do(func() {
		instance = &Config{}
		err = instance.loadWithRootPath(rootPath)
	})
	return err
}

func InitWithContent(content []byte) error {
	var err error
	once.Do(func() {
		instance = &Config{}
		err = instance.loadWithContent(content)
	})
	return err
}

type Config struct {
	SMSGateway    SMSGateway    `yaml:"smsGateway"`
	CustomSenders []Sender      `yaml:"customSenders"`
	Preferences   []Preferences `yaml:"preferences"`
	Weather       API           `yaml:"weather"`
	Cyclocity     API           `yaml:"cyclocity"`
	Navitia       API           `yaml:"navitia"`
}

func (config *Config) loadWithRootPath(rootPath string) error {
	absPath, err := filepath.Abs(rootPath + "/config.yml")
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(absPath)
	return config.loadWithContent(content)
}

func (config *Config) loadWithContent(content []byte) error {
	content = []byte(os.ExpandEnv(string(content))) // expand with env variables first
	err := yaml.Unmarshal(content, config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) FindCustomSender(phoneNumber string) *Sender {
	for _, sender := range config.CustomSenders {
		if sender.PhoneNumber == phoneNumber {
			return &sender
		}
	}
	return nil
}

func (config *Config) GetPreferences(phoneNumber string) *Preferences {
	for _, preferences := range config.Preferences {
		if preferences.PhoneNumber == phoneNumber {
			return &preferences
		}
	}
	return nil
}

func (config *Config) GetDefaults(phoneNumber string) *Defaults {
	preferences := config.GetPreferences(phoneNumber)
	if preferences != nil {
		return &preferences.Defaults
	}
	return nil
}
