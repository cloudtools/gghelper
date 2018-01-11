package main

import (
	"flag"
	"fmt"
	"os"
)

// CreateGroup - command execution for creategroup subcommand
func CreateGroup(args []string) {

	var region, profile string
	var name string

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.StringVar(&name, "name", "", "core name")
	cmdFlags.Parse(args)

	if name == "" {
		fmt.Printf("Must specify a name for the core\n")
		os.Exit(1)
	}

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.CreateGroup(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	ggSession.GetConfigArchive()
	ggSession.WriteGroupConfig("config.json")
}
