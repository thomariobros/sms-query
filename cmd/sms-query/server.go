package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"

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
		processQuery(w, r, parseParamsNexmo)
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
	w.WriteHeader(http.StatusOK)
}

type paramsParser func(request *http.Request) (string, string)

func parseParamsNexmo(request *http.Request) (string, string) {
	var from, query string
	if request.Method == "GET" {
		from = strings.TrimSpace(request.URL.Query().Get("msisdn"))
		query = strings.TrimSpace(request.URL.Query().Get("text"))
	} else if request.Method == "POST" {
		request.ParseForm()
		from = strings.TrimSpace(request.Form.Get("msisdn"))
		query = strings.TrimSpace(request.Form.Get("text"))
	}
	if from != "" && query != "" {
		return from, query
	}
	return "", ""
}
