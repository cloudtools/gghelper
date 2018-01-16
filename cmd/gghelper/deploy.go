package main

import (
	"flag"
	"fmt"
)

// CreateDeployment - command execution for createdeployment subcommand
func CreateDeployment(args []string) {

	var region, profile string

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.Parse(args)

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.CreateDeployment()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// ListDeployment - command execution for listdeployment subcommand
func ListDeployment(args []string) {

	var region, profile string

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.Parse(args)

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.ListDeployment()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// ResetDeployment - command execution for resetdeployment subcommand
func ResetDeployment(args []string) {

	var region, profile string
	var force bool

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.BoolVar(&force, "f", false, "force reset")
	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.Parse(args)

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.ResetDeployment(force)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
