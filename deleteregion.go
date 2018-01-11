package gghelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/aws/aws-sdk-go/service/iot"
)

// DeleteRegion - delete iot/greengrass objects from region
func (ggSession *GreengrassSession) DeleteRegion() error {
	things, err := ggSession.iot.ListThings(&iot.ListThingsInput{})
	if err != nil {
		fmt.Printf("Err: ListThings - %v\n", err)
		return err
	}
	// fmt.Printf("things: %v\n", things)

	for _, t := range things.Things {
		thingName := *t.ThingName
		// fmt.Printf("ThingName %s\n", *t.ThingName)
		principals, err := ggSession.iot.ListThingPrincipals(&iot.ListThingPrincipalsInput{
			ThingName: &thingName,
		})
		if err != nil {
			fmt.Printf("Error ListThingPrincipals - %v\n", err)
			return err
		}
		for _, p := range principals.Principals {
			// fmt.Printf("principal %s\n", *p)
			_, err = ggSession.iot.DetachThingPrincipal(&iot.DetachThingPrincipalInput{
				ThingName: &thingName,
				Principal: p,
			})
			if err != nil {
				return err
			}
			fmt.Printf("Deleted principal %s from thing %s\n", *p, thingName)
		}
		_, err = ggSession.iot.DeleteThing(&iot.DeleteThingInput{
			ThingName: &thingName,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Deleted thing %s\n", thingName)
	}

	// Delete policies
	policies, err := ggSession.iot.ListPolicies(&iot.ListPoliciesInput{})
	for _, p := range policies.Policies {
		// fmt.Printf("Policy %s\n", *p)
		principals, _ := ggSession.iot.ListPolicyPrincipals(&iot.ListPolicyPrincipalsInput{
			PolicyName: p.PolicyName,
		})
		for _, principal := range principals.Principals {
			ggSession.iot.DetachPrincipalPolicy(&iot.DetachPrincipalPolicyInput{
				PolicyName: p.PolicyName,
				Principal:  principal,
			})
			fmt.Printf("Detached policy principal %s from %s\n", *principal, *p.PolicyName)
		}
		_, err = ggSession.iot.DeletePolicy(&iot.DeletePolicyInput{
			PolicyName: p.PolicyName,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Deleted policy %s\n", *p.PolicyName)
	}

	// Delete certificates
	certificates, err := ggSession.iot.ListCertificates(&iot.ListCertificatesInput{})
	if err != nil {
		fmt.Printf("Err: ListCertificates - %v\n", err)
		return err
	}

	inactive := "INACTIVE"
	for _, cert := range certificates.Certificates {
		ggSession.iot.UpdateCertificate(&iot.UpdateCertificateInput{
			CertificateId: cert.CertificateId,
			NewStatus:     &inactive,
		})
		ggSession.iot.DeleteCertificate(&iot.DeleteCertificateInput{
			CertificateId: cert.CertificateId,
		})
		fmt.Printf("Deleting certificate: %s\n", *cert.CertificateId)
	}

	// Delete groups
	groups, err := ggSession.greengrass.ListGroups(&greengrass.ListGroupsInput{})
	if err != nil {
		fmt.Printf("Err: ListGroups - %v\n", err)
		return err
	}
	for _, g := range groups.Groups {
		ggSession.greengrass.ResetDeployments(&greengrass.ResetDeploymentsInput{
			GroupId: g.Id,
		})
		_, err = ggSession.greengrass.DeleteGroup(&greengrass.DeleteGroupInput{
			GroupId: g.Id,
		})
		if err != nil {
			fmt.Printf("Err: DeleteGroup - %v\n", err)
			return err
		}
		fmt.Printf("Deleted group: %s\n", *g.Id)
	}

	// Delete core definitions
	cores, err := ggSession.greengrass.ListCoreDefinitions(&greengrass.ListCoreDefinitionsInput{})
	if err != nil {
		fmt.Printf("Err: ListCoreDefinitions - %v\n", err)
		return err
	}
	for _, core := range cores.Definitions {
		ggSession.greengrass.DeleteCoreDefinition(&greengrass.DeleteCoreDefinitionInput{
			CoreDefinitionId: core.Id,
		})
		if err != nil {
			fmt.Printf("Err: DeleteCoreDefinition - %v\n", err)
			return err
		}
		fmt.Printf("Deleted core definition: %s\n", *core.Id)
	}

	// Delete subscriptions
	subs, err := ggSession.greengrass.ListSubscriptionDefinitions(&greengrass.ListSubscriptionDefinitionsInput{})
	for _, sub := range subs.Definitions {
		ggSession.greengrass.DeleteSubscriptionDefinition(&greengrass.DeleteSubscriptionDefinitionInput{
			SubscriptionDefinitionId: sub.Id,
		})
		if err != nil {
			fmt.Printf("Err: DeleteSubscriptionDefinition - %v\n", err)
			return err
		}
		fmt.Printf("Deleted subscription definition: %s\n", *sub.Id)
	}

	// Delete functions
	functions, err := ggSession.greengrass.ListFunctionDefinitions(&greengrass.ListFunctionDefinitionsInput{})
	for _, function := range functions.Definitions {
		ggSession.greengrass.DeleteFunctionDefinition(&greengrass.DeleteFunctionDefinitionInput{
			FunctionDefinitionId: function.Id,
		})
		if err != nil {
			fmt.Printf("Err: DeleteFunctionDefinition - %v\n", err)
			return err
		}
		fmt.Printf("Deleted function definition: %s\n", *function.Id)
	}

	// Delete loggers
	loggers, err := ggSession.greengrass.ListLoggerDefinitions(&greengrass.ListLoggerDefinitionsInput{})
	for _, logger := range loggers.Definitions {
		ggSession.greengrass.DeleteLoggerDefinition(&greengrass.DeleteLoggerDefinitionInput{
			LoggerDefinitionId: logger.Id,
		})
		if err != nil {
			fmt.Printf("Err: DeleteLoggerDefinition - %v\n", err)
			return err
		}
		fmt.Printf("Deleted logger definition: %s\n", *logger.Id)
	}

	ggSession.config.ClearConfig()

	return nil
}
