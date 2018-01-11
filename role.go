package gghelper

import (
	"github.com/aws/aws-sdk-go/service/iam"
)

// RoleLookup - lookup a Role name and return the associated ARN
func (ggSession *GreengrassSession) RoleLookup(role string) (*string, error) {
	getRoleInput := iam.GetRoleInput{
		RoleName: &role,
	}
	getRoleOutput, err := ggSession.iam.GetRole(&getRoleInput)
	if err != nil {
		return nil, err
	}
	return getRoleOutput.Role.Arn, nil

}
