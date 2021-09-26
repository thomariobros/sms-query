package run

import (
	"errors"
	"os"

	"github.com/golang/glog"

	"main/pkg/command"
	"main/pkg/config"
	"main/pkg/i18n"
)

func Init(content map[string][]byte) {
	// load config
	err := config.InitWithContent(content["config"])
	if err != nil {
		glog.Fatalf("error while loading config: %s", err)
		os.Exit(1)
	}
	// load i18n
	i18n.LoadWithContent(content)
}

func ExecuteQuery(from, query string) (string, error) {
	if !checkInputs(from, query) {
		return "", errors.New("wrong inputs")
	}
	return command.NewDispatcher().Execute(from, query), nil
}

func ExecuteQuerySend(from, query string, background bool) error {
	if !checkInputs(from, query) {
		return errors.New("wrong inputs")
	}
	if background {
		// new thread to execute job in background
		go command.NewDispatcher().ExecuteSend(from, query)
		return nil
	}
	return command.NewDispatcher().ExecuteSend(from, query)
}

func checkInputs(from string, query string) bool {
	if len(from) == 0 || len(from) > 15 || len(query) == 0 && len(query) > 100 {
		return false
	}
	// check if phone number is allowed
	if !config.GetInstance().SMSGateway.IsPhoneNumberAllowed(from) {
		glog.Warningf("phone number %s not allowed", from)
		return false
	}
	return true
}
