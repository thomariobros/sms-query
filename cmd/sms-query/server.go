package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/glog"

	"main/pkg/config"
	"main/pkg/run"
)

// Version build version
var Version = "latest"

func main() {
	// cli params
	bind := flag.String("bind", "127.0.0.1:8080", "bind address")
	flag.Parse()

	// init
	assets := map[string]string{
		"config":     "config/config.yml",
		"i18n-en_US": "config/i18n/en_US.json",
		"i18n-fr_FR": "config/i18n/fr_FR.json",
	}
	content := make(map[string][]byte)
	for key, path := range assets {
		asset, err := Asset(path)
		if err != nil {
			glog.Fatalf("Error while loading config: %s", err)
			os.Exit(1)
		}
		content[key] = asset
	}
	run.Init(content)

	// handlers
	registerHandlers()

	glog.Info(fmt.Sprintf("sms-query %s - http server listens and serves on %v...", Version, *bind))
	err := http.ListenAndServe(*bind, nil)
	if err != nil {
		glog.Fatalf("ListenAndServe: %s", err)
		os.Exit(1)
	}
}

func registerHandlers() {
	// query from nexmo
	http.HandleFunc("/nexmo", func(w http.ResponseWriter, r *http.Request) {
		processQuery(w, r, parseNexmo)
	})
}

func processQuery(w http.ResponseWriter, r *http.Request, parserFn paramsParser) {
	from, query := parserFn(r)
	err := run.ExecuteQuerySend(from, query, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// immediate http response
	w.WriteHeader(http.StatusNoContent)
}

type paramsParser func(request *http.Request) (string, string)

func parseNexmo(request *http.Request) (string, string) {
	// check jwt signature https://developer.nexmo.com/messages/concepts/signed-webhooks
	toString := strings.Split(request.Header.Get("authorization"), " ")[1]
	token, err := jwt.Parse(toString, func(token *jwt.Token) (interface{}, error) {
		// validate that alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected jwt signing method: %v", token.Header["alg"])
		}
		signatureSecret := config.GetInstance().SMSGateway.Provider.SignatureSecret
		return []byte(signatureSecret), nil
	})
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		glog.Errorf("Error while checking jwt signature: %s", err)
		return "", ""
	}
	// get from and query params
	var from, query string
	if request.Method == "POST" {
		var message nexmoMessage
		if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
			glog.Errorf("Error while decoding JSON: %s", err)
			return "", ""
		}
		from = strings.TrimSpace(message.From.Number)
		query = strings.TrimSpace(message.Message.Content.Text)
	}
	if from != "" && query != "" {
		return from, query
	}
	return "", ""
}

type nexmoMessage struct {
	From struct {
		Type   string `json:"type"`
		Number string `json:"number"`
	} `json:"from"`
	Message struct {
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}
