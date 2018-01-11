package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

// CreateCore - command execution for createcore subcommand
func CreateCore(args []string) {

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

	_, err := ggSession.CreateCore(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	buf := new(bytes.Buffer)
	ggSession.WriteGGConfig(buf)
	fmt.Println(buf.String())

	ggSession.GetConfigArchive()

	//ggSession.Cleanup()
}
