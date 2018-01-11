package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// DeleteRegion - command execution for deleteregion subcommand
func DeleteRegion(args []string) {

	var region, profile string
	var yes bool

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.BoolVar(&yes, "y", false, "yes - remove iot/greengrass objects from the entire region")
	cmdFlags.Parse(args)

	if !yes {
		stdinReader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter 'yes' to remove all iot/greengrass objects from the region %s: ", region)
		text, _ := stdinReader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		if text != "yes" {
			fmt.Printf("Must agree to removing region objects by typing 'yes'\n")
			os.Exit(1)
		}
	}

	ggSession := CreateSession(profile, region)

	ggSession.DeleteRegion()

	ggSession.WriteGroupConfig("config.json")
}
