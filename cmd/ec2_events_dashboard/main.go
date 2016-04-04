package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	dashboard "github.com/jrab89/ec2_events_dashboard"

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

	ec2Clients, err := dashboard.NewClientsFromCreds(opts.Creds)
	if err != nil {
		// TODO: print out some error message
		os.Exit(1)
	}

	instances, err := dashboard.InstancesWithEvents(ec2Clients...)
	if err != nil {
		// TODO: print out some error message
		os.Exit(1)
	}

	fmt.Println(len(instances))
}
