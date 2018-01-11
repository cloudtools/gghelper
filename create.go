package gghelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/aws/aws-sdk-go/service/iot"
)

// CreateThing - create a new Thing object
func (ggSession *GreengrassSession) CreateThing(name string) (*iot.CreateThingOutput, error) {
	var err error

	// Create a new set of keys and certificate
	setAsActive := true
	ggSession.keyCertOutput, err = ggSession.iot.CreateKeysAndCertificate(&iot.CreateKeysAndCertificateInput{SetAsActive: &setAsActive})
	//fmt.Printf("keyCertOutput: %+v\n", keyOutput)
	if err != nil {
		return nil, err
	}
	fmt.Printf("CertificateId: %s\n", *ggSession.keyCertOutput.CertificateId)

	// Create a "thing"
	thingOutput, err := ggSession.iot.CreateThing(&iot.CreateThingInput{
		ThingName: &name,
	})
	if err != nil {
		fmt.Printf("CreateThing error: %v\n", err)
		return nil, err
	}
	fmt.Printf("ThingArn: %s\n", *thingOutput.ThingArn)

	// Attach the thing principal
	_, err = ggSession.iot.AttachThingPrincipal(&iot.AttachThingPrincipalInput{
		Principal: ggSession.keyCertOutput.CertificateArn,
		ThingName: &name,
	})
	if err != nil {
		fmt.Printf("AttachThingPrincipal error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Called AttachThingPrincipal policy\n")

	return thingOutput, nil
}

// CreateThingPolicy - create the policy for a thing and attach it
func (ggSession *GreengrassSession) CreateThingPolicy(name string) error {
	// Get or create the IoT policy
	policyName := fmt.Sprintf("%s-policy", name)
	_, err := ggSession.iot.GetPolicy(&iot.GetPolicyInput{
		PolicyName: &policyName,
	})
	if err == nil {
		fmt.Printf("Found existing policy: %s\n", policyName)
	} else {
		_, err = ggSession.iot.CreatePolicy(&iot.CreatePolicyInput{
			PolicyName:     &policyName,
			PolicyDocument: &policyDocument,
		})
		if err != nil {
			fmt.Printf("CreatePolicy error: %v\n", err)
			return err
		}
		fmt.Printf("Created policy: %s\n", policyName)
	}

	// Attach the principal policy
	_, err = ggSession.iot.AttachPrincipalPolicy(&iot.AttachPrincipalPolicyInput{
		PolicyName: &policyName,
		Principal:  ggSession.keyCertOutput.CertificateArn,
	})
	if err != nil {
		fmt.Printf("AttachPrincipalPolicy error: %v\n", err)
		return err
	}
	fmt.Printf("Called AttachPrincipalPolicy\n")
	return nil
}

// CreateCore - create a new Greengrass Core object
func (ggSession *GreengrassSession) CreateCore(thing string) (*iot.CreateThingOutput, error) {
	thingOutput, err := ggSession.CreateThing(thing)
	if err != nil {
		return nil, err
	}

	err = ggSession.CreateThingPolicy(thing)
	if err != nil {
		return nil, err
	}

	// Update the configuration
	certID := (*ggSession.keyCertOutput.CertificateId)[0:10]
	ggSession.ggconfig.CoreThing.CertPath = fmt.Sprintf("%s.cert.pem", certID)
	ggSession.ggconfig.CoreThing.KeyPath = fmt.Sprintf("%s.private.key", certID)
	ggSession.ggconfig.CoreThing.CAPath = "root.ca.pem"
	ggSession.ggconfig.CoreThing.GGHost = fmt.Sprintf("greengrass.iot.%s.amazonaws.com", *ggSession.session.Config.Region)

	endpoint, _ := ggSession.iot.DescribeEndpoint(&iot.DescribeEndpointInput{})
	ggSession.ggconfig.CoreThing.IOTHost = *endpoint.EndpointAddress
	ggSession.ggconfig.Runtime.Cgroup.UseSystemd = "yes"

	// Update configuration data
	ggSession.config.Core.ThingName = thing
	ggSession.config.Core.ThingArn = *thingOutput.ThingArn
	ggSession.config.Core.CertID = *ggSession.keyCertOutput.CertificateId
	ggSession.config.Core.CertArn = *ggSession.keyCertOutput.CertificateArn

	ggSession.ggconfig.CoreThing.ThingArn = *thingOutput.ThingArn

	return thingOutput, nil
}

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
