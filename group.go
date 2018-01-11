package gghelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/greengrass"
)

// CreateGroup - create a new Greengrass group
func (ggSession *GreengrassSession) CreateGroup(name string) error {
	// fmt.Printf("creategroup: %v\n", groupOutput)
	thingOutput, err := ggSession.CreateCore(name)
	if err != nil {
		fmt.Printf("CreateCore error: %v\n", err)
	}
	fmt.Printf("Created core '%s'\n", name)

	coreDefinition, err := ggSession.greengrass.CreateCoreDefinition(&greengrass.CreateCoreDefinitionInput{
		Name: &name,
	})
	if err != nil {
		fmt.Printf("CreateCoreDefinition error: %v\n", err)
		return err
	}
	fmt.Printf("Created core definition\n")

	definitionInput := greengrass.CreateCoreDefinitionVersionInput{
		CoreDefinitionId: coreDefinition.Id,
		Cores: []*greengrass.Core{
			&greengrass.Core{
				CertificateArn: ggSession.keyCertOutput.CertificateArn,
				Id:             coreDefinition.Id,
				ThingArn:       thingOutput.ThingArn,
			},
		},
	}

	definitionVersion, err := ggSession.greengrass.CreateCoreDefinitionVersion(&definitionInput)
	if err != nil {
		fmt.Printf("CreateCoreDefinitionVersion error: %v\n", err)
		return err
	}
	fmt.Printf("Created core definition version\n")

	groupOutput, err := ggSession.greengrass.CreateGroup(&greengrass.CreateGroupInput{Name: &name})
	if err != nil {
		fmt.Printf("creategroup error: %v\n", err)
		return err
	}
	fmt.Printf("Created group\n")

	// Update group configuration
	ggSession.config.CoreDefinition.ID = *definitionVersion.Id
	ggSession.config.CoreDefinition.VersionArn = *definitionVersion.Arn
	ggSession.config.Group.ID = *groupOutput.Id

	ggSession.updateGroup()

	return nil
}

func (ggSession *GreengrassSession) updateGroup() error {
	input := &greengrass.CreateGroupVersionInput{
		GroupId: &ggSession.config.Group.ID,
	}

	if ggSession.config.CoreDefinition.VersionArn != "" {
		input.CoreDefinitionVersionArn = &ggSession.config.CoreDefinition.VersionArn
	}
	if ggSession.config.DeviceDefinition.VersionArn != "" {
		input.DeviceDefinitionVersionArn = &ggSession.config.DeviceDefinition.VersionArn
	}
	if ggSession.config.FunctionDefinition.VersionArn != "" {
		input.FunctionDefinitionVersionArn = &ggSession.config.FunctionDefinition.VersionArn
	}
	if ggSession.config.LoggerDefinition.VersionArn != "" {
		input.LoggerDefinitionVersionArn = &ggSession.config.LoggerDefinition.VersionArn
	}
	if ggSession.config.SubscriptionDefinition.VersionArn != "" {
		input.SubscriptionDefinitionVersionArn = &ggSession.config.SubscriptionDefinition.VersionArn
	}

	_, err := ggSession.greengrass.CreateGroupVersion(input)
	if err != nil {
		fmt.Printf("updateGroup error: %v\n", err)
		return err
	}
	fmt.Printf("Updated group version\n")

	return nil
}

// ListGroup - list the group definition
func (ggSession *GreengrassSession) ListGroup() error {
	group, err := ggSession.greengrass.GetGroup(&greengrass.GetGroupInput{
		GroupId: &ggSession.config.Group.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("group: %v\n", group)

	groupVersion, err := ggSession.greengrass.GetGroupVersion(&greengrass.GetGroupVersionInput{
		GroupId:        group.Id,
		GroupVersionId: group.LatestVersion,
	})
	if err != nil {
		return err
	}
	fmt.Printf("group version: %v\n", groupVersion)

	return nil
}
