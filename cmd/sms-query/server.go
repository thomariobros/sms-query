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

	"sms-query/pkg/config"
	"sms-query/pkg/run"
)

// Version or commit id as version
var Version = "latest"
var CommitId = os.Getenv("COMMIT_ID")

func main() {
	// cli params
	config := flag.String("config", "config/config.yml", "config path")
	bind := flag.String("bind", "0.0.0.0:8080", "bind address")
	flag.Parse()
	// force logtostderr
	flag.Lookup("logtostderr").Value.Set("true")

	// init
	run.Init(*config)

	// handlers
	registerHandlers()

	version := Version
	if version == "latest" && CommitId != "" {
		version = CommitId
	}
	glog.Info(fmt.Sprintf("sms-query %s - http server listens and serves on %v...", version, *bind))
	err := http.ListenAndServe(*bind, nil)
	if err != nil {
		glog.Fatalf("ListenAndServe: %s", err)
		os.Exit(1)
	}
}

func registerHandlers() {
	http.HandleFunc("/nexmo/inbound", func(w http.ResponseWriter, r *http.Request) {
		processInboundQuery(w, r, parseNexmoInbound)
	})
	http.HandleFunc("/nexmo/status", func(w http.ResponseWriter, r *http.Request) {
		processStatus(w, r, parseNexmoStatus)
	})
}

func processInboundQuery(w http.ResponseWriter, r *http.Request, parserFn inboundQueryParser) {
	from, query := parserFn(r)
	err := run.ExecuteQuerySend(from, query, true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// immediate http response
	w.WriteHeader(http.StatusNoContent)
}

type inboundQueryParser func(request *http.Request) (string, string)

func parseNexmoInbound(request *http.Request) (string, string) {
	ok := checkSignature(request)
	if !ok {
		return "", ""
	}
	// get from and query params
	var from, query string
	if request.Method == "POST" {
		var message nexmoInboundMessage
		if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
			glog.Errorf("Error while decoding JSON: %s", err)
			return "", ""
		}
		from = message.From.Number
		query = message.Message.Content.Text
	}
	return from, query
}

type nexmoInboundMessage struct {
	From struct {
		Number string `json:"number"`
	} `json:"from"`
	Message struct {
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`
}

func processStatus(w http.ResponseWriter, r *http.Request, parserFn statusParser) {
	to, status := parserFn(r)
	glog.Infof("Status to=%s, status=", to, status)
	w.WriteHeader(http.StatusNoContent)
}

type statusParser func(request *http.Request) (string, string)

func parseNexmoStatus(request *http.Request) (string, string) {
	ok := checkSignature(request)
	if !ok {
		return "", ""
	}
	var to, status string
	if request.Method == "POST" {
		var message nexmoStatusMessage
		if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
			glog.Errorf("Error while decoding JSON: %s", err)
			return "", ""
		}
	}
	return to, status
}

type nexmoStatusMessage struct {
	To struct {
		Number string `json:"number"`
	} `json:"to"`
	Status string `json:"status"`
}

func checkSignature(request *http.Request) bool {
	// check jwt signature https://developer.nexmo.com/messages/concepts/signed-webhooks
	authorizationHeader := request.Header.Get("authorization")
	if authorizationHeader == "" {
		return false
	}
	splitAuthorizationHeader := strings.Split(authorizationHeader, " ")
	if len(splitAuthorizationHeader) != 2 {
		return false
	}
	tokenString := splitAuthorizationHeader[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate that alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected jwt signing method: %v", token.Header["alg"])
		}
		signatureSecret := config.GetInstance().SMSGateway.Provider.SignatureSecret
		return []byte(signatureSecret), nil
	})
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		glog.Errorf("Error while checking jwt signature: %s", err)
		return false
	}
	return token.Valid
}
