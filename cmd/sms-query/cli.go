package main

import (
	"flag"
	"fmt"
	"os"

	"main/pkg/run"

	"github.com/golang/glog"
)

// Version build version
var Version = "latest"

func main() {
	// cli params
	from := flag.String("from", "toto", "from phone number")
	query := flag.String("query", "help", "query")
	send := flag.Bool("send", false, "send")
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

	var output string
	var err error
	if *send {
		err = run.ExecuteQuerySend(*from, *query, false)
		if err != nil {
			output = err.Error()
		}
	} else {
		output, err = run.ExecuteQuery(*from, *query)
		if err != nil {
			output = err.Error()
		}
		if output == "" {
			output = "no output"
		}
	}
	if output != "" {
		fmt.Println(output)
	}
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
