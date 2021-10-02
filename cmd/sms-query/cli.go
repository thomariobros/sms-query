package main

import (
	"flag"
	"fmt"
	"os"

	"sms-query/pkg/run"
)

// Version build version
var Version = "latest"

func main() {
	// cli params
	config := flag.String("config", "config/config.yml", "config path")
	from := flag.String("from", "toto", "from phone number")
	query := flag.String("query", "help", "query")
	send := flag.Bool("send", false, "send")
	flag.Parse()

	// init
	run.Init(*config)

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
