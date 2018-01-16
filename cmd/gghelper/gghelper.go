package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cloudtools/gghelper"
)

func printUsage() {
	fmt.Println("usage: gghelper <command> [<args>]")
	fmt.Println("Commands: ")
	fmt.Println(" createcore       - create a new Greengrass core")
	fmt.Println(" createdeployment - create a new Greengrass core")
	fmt.Println(" creategroup      - create a new Greengrass group")
	fmt.Println(" createsub        - create a new Greengrass subscription")
	fmt.Println(" deleteregion     - delete all greengrass/iot resources in a region")
	fmt.Println(" lambda           - submit a lambda function")
	fmt.Println(" listdeployment   - list greengrass deployment")
	fmt.Println(" listgroup        - list greengrass group")
	fmt.Println(" listlambda       - list greengrass view of lambda functions")
	fmt.Println(" listsub          - list greengrass subscriptions")
	fmt.Println(" resetdeployment  - reset greengrass deployments")
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "createcore":
		CreateCore(os.Args[2:])
	case "createdeployment", "createdeploy", "deploy":
		CreateDeployment(os.Args[2:])
	case "creategroup":
		CreateGroup(os.Args[2:])
	case "createsub":
		CreateSub(os.Args[2:])
	case "deleteregion":
		DeleteRegion(os.Args[2:])
	case "lambda":
		Lambda(os.Args[2:])
	case "listdeployment", "listdeploy":
		ListDeployment(os.Args[2:])
	case "listgroup":
		ListGroup(os.Args[2:])
	case "listlambda":
		ListLambda(os.Args[2:])
	case "listsub":
		ListSub(os.Args[2:])
	case "resetdeployment", "resetdeploy":
		ResetDeployment(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

// CreateSession - common function for creating an AWS session
func CreateSession(profile, region string) *gghelper.GreengrassSession {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if region != "" {
		options.Config = aws.Config{Region: aws.String(region)}
	}
	if profile != "" {
		options.Profile = profile
	}

	session := session.Must(session.NewSessionWithOptions(options))

	return gghelper.NewGreengrassSession(session)
}
