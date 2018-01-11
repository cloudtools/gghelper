package gghelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/satori/go.uuid"
)

func (ggSession *GreengrassSession) mapSubToArn(name string) string {
	if name == "cloud" {
		return name
	}

	// Check if name is a lambda function
	if val, ok := ggSession.config.LambdaFunctions[name]; ok {
		return val.Arn
	}

	// Todo: Check if name is a device

	return ""
}

// CreateSub - create a greengrass subscription
func (ggSession *GreengrassSession) CreateSub(source, target, subject string) error {
	sourceArn := ggSession.mapSubToArn(source)
	targetArn := ggSession.mapSubToArn(target)

	newUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	uuidString := newUUID.String()

	// Check if we need to create the initial version
	if ggSession.config.SubscriptionDefinition.ID == "" {
		newSubscription, err := ggSession.greengrass.CreateSubscriptionDefinition(&greengrass.CreateSubscriptionDefinitionInput{
			InitialVersion: &greengrass.SubscriptionDefinitionVersion{
				Subscriptions: []*greengrass.Subscription{
					&greengrass.Subscription{
						Source:  &sourceArn,
						Target:  &targetArn,
						Subject: &subject,
						Id:      &uuidString,
					},
				},
			},
		})

		if err != nil {
			return err
		}
		fmt.Printf("Created new subscription\n")
		ggSession.config.SubscriptionDefinition.ID = *newSubscription.Id
		ggSession.config.SubscriptionDefinition.VersionArn = *newSubscription.LatestVersionArn

		ggSession.updateGroup()

		return nil
	}

	// Add subscription to existing
	subscription, _ := ggSession.greengrass.GetSubscriptionDefinition(&greengrass.GetSubscriptionDefinitionInput{
		SubscriptionDefinitionId: &ggSession.config.SubscriptionDefinition.ID,
	})

	subscriptionVersion, _ := ggSession.greengrass.GetSubscriptionDefinitionVersion(&greengrass.GetSubscriptionDefinitionVersionInput{
		SubscriptionDefinitionId:        subscription.Id,
		SubscriptionDefinitionVersionId: subscription.LatestVersion,
	})
	subscriptions := subscriptionVersion.Definition.Subscriptions

	subscriptions = append(subscriptions, &greengrass.Subscription{
		Source:  &sourceArn,
		Target:  &targetArn,
		Subject: &subject,
		Id:      &uuidString,
	})

	output, err := ggSession.greengrass.CreateSubscriptionDefinitionVersion(&greengrass.CreateSubscriptionDefinitionVersionInput{
		SubscriptionDefinitionId: subscription.Id,
		Subscriptions:            subscriptions,
	})
	if err != nil {
		return err
	}

	ggSession.config.SubscriptionDefinition.VersionArn = *output.Arn
	fmt.Printf("Updated subscription\n")

	ggSession.updateGroup()

	return nil
}

// ListSub - list a greengrass subscription
func (ggSession *GreengrassSession) ListSub() error {
	if ggSession.config.SubscriptionDefinition.ID == "" {
		fmt.Printf("No initial subscription defined\n")
		return nil
	}
	subscription, err := ggSession.greengrass.GetSubscriptionDefinition(&greengrass.GetSubscriptionDefinitionInput{
		SubscriptionDefinitionId: &ggSession.config.SubscriptionDefinition.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("subscription: %v\n", subscription)
	subscriptionVersion, err := ggSession.greengrass.GetSubscriptionDefinitionVersion(&greengrass.GetSubscriptionDefinitionVersionInput{
		SubscriptionDefinitionId:        subscription.Id,
		SubscriptionDefinitionVersionId: subscription.LatestVersion,
	})
	if err != nil {
		return err
	}
	fmt.Printf("subscription version: %v\n", subscriptionVersion)

	return nil
}
