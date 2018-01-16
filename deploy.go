package gghelper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/greengrass"
)

// CreateDeployment - Create a deployment
func (ggSession *GreengrassSession) CreateDeployment() error {
	deployments, err := ggSession.greengrass.ListDeployments(&greengrass.ListDeploymentsInput{
		GroupId: &ggSession.config.Group.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("deployments: %v\n", deployments)

	newDeployment := "NewDeployment"
	create, err := ggSession.greengrass.CreateDeployment(&greengrass.CreateDeploymentInput{
		DeploymentType: &newDeployment,
		GroupId:        &ggSession.config.Group.ID,
		GroupVersionId: &ggSession.config.Group.Version,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", create)

	return nil
}

// ListDeployment - lists the current deployment status
func (ggSession *GreengrassSession) ListDeployment() error {
	deployments, err := ggSession.greengrass.ListDeployments(&greengrass.ListDeploymentsInput{
		GroupId: &ggSession.config.Group.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("deployments: %v\n", deployments)

	for _, d := range deployments.Deployments {
		deployment, err := ggSession.greengrass.GetDeploymentStatus(&greengrass.GetDeploymentStatusInput{
			DeploymentId: d.DeploymentId,
			GroupId:      &ggSession.config.Group.ID,
		})
		if err != nil {
			return err
		}
		fmt.Printf("deployment: %v\n", deployment)
	}

	return nil
}

// ResetDeployment - reset deployments status
func (ggSession *GreengrassSession) ResetDeployment(force bool) error {
	reset, err := ggSession.greengrass.ResetDeployments(&greengrass.ResetDeploymentsInput{
		GroupId: &ggSession.config.Group.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("deployment reset: %v\n", reset)

	return nil
}
