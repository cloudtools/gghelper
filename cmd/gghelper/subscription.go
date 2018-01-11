package main

import (
	"flag"
	"fmt"
	"os"
)

// CreateSub - command execution for createsub subcommand
func CreateSub(args []string) {

	var region, profile string
	var source, subject, target string

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.StringVar(&source, "source", "", "source name")
	cmdFlags.StringVar(&target, "target", "", "target name")
	cmdFlags.StringVar(&subject, "subject", "", "subject name")
	cmdFlags.Parse(args)

	if source == "" {
		fmt.Printf("Must specify a source for the subscription\n")
		os.Exit(1)
	}
	if target == "" {
		fmt.Printf("Must specify a source for the subscription\n")
		os.Exit(1)
	}

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.CreateSub(source, target, subject)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	ggSession.WriteGroupConfig("config.json")
}

// ListSub - command execution for listsub subcommand
func ListSub(args []string) {

	var region, profile string

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.Parse(args)

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.ListSub()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
