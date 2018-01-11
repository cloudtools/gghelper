package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudtools/gghelper"
)

// Lambda - command execution for lambda subcommand
func Lambda(args []string) {

	var region, profile string
	var directory, handler, role, runtime, name string
	var pinned bool

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.StringVar(&directory, "d", "", "code directory")
	cmdFlags.StringVar(&handler, "handler", "", "handler")
	cmdFlags.StringVar(&name, "name", "", "lambda name")
	cmdFlags.BoolVar(&pinned, "pinned", false, "pinned - long lived function")
	cmdFlags.StringVar(&role, "role", "", "role")
	cmdFlags.StringVar(&runtime, "runtime", "python2.7", "runtime")
	cmdFlags.Parse(args)

	if directory == "" {
		fmt.Printf("Must specify code directory\n")
		os.Exit(1)
	}

	if handler == "" {
		fmt.Printf("Must specify handler\n")
		os.Exit(1)
	}

	if name == "" {
		fmt.Printf("Must specify lambda name\n")
		os.Exit(1)
	}

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	roleArn, err := ggSession.RoleLookup(role)
	if err != nil {
		fmt.Printf("Could not find arn for role %s\n", role)
		os.Exit(1)
	}

	config := gghelper.LambdaConfig{
		Directory: directory,
		Handler:   handler,
		Name:      name,
		Pinned:    pinned,
		RoleArn:   *roleArn,
		Runtime:   runtime,
	}
	err = ggSession.SubmitFunction(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	ggSession.WriteGroupConfig("config.json")
}

// ListLambda - command execution for listlambda subcommand
func ListLambda(args []string) {

	var region, profile string

	cmdFlags := flag.NewFlagSet("ask", flag.ExitOnError)

	cmdFlags.StringVar(&region, "r", "", "region")
	cmdFlags.StringVar(&profile, "p", "", "profile")
	cmdFlags.Parse(args)

	ggSession := CreateSession(profile, region)

	ggSession.LoadGroupConfig("config.json")

	err := ggSession.ListLambda()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
