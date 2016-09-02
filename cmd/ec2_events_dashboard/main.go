package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	ec2_events_dashboard "../.."
	"github.com/jessevdk/go-flags"

	"os"
)

func main() {
	var opts struct {
		Creds []string `long:"creds" short:"c" description:"AWS API keys (this flag can be used more than once) in the form of <YOUR_ACCESS_KEY_ID>:<YOUR_SECRET_ACCESS_KEY>"`
	}

	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	ec2Clients, err := ec2_events_dashboard.NewClientsFromCreds(opts.Creds)
	if err != nil {
		// TODO: print out some error message
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		instances, err := ec2_events_dashboard.InstancesWithEvents(ec2Clients...)
		if err != nil {
			// TODO: print out some error message
			os.Exit(1)
		}

		// TODO: probably should 500 if this errors?
		marshalled, _ := json.MarshalIndent(instances, "", "    ")
		fmt.Fprintf(w, string(marshalled))
	})
	http.ListenAndServe(":3000", nil)
}
