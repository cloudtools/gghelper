package gghelper

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var policyDocument = `{
	"Version": "2012-10-17",
	"Statement": [
	  {
		"Effect": "Allow",
		"Action": [
		  "iot:*",
		  "greengrass:*"
		],
		"Resource": [
		  "*"
		]
	  }
	]
  }`

type ggHelperConfig struct {
	Core struct {
		CertArn   string `json:"cert_arn"`
		CertID    string `json:"cert_id"`
		ThingArn  string `json:"thing_arn"`
		ThingName string `json:"thing_name"`
	}
}

// GreengrassSession is an object for Greengrass sessions
type GreengrassSession struct {
	session       *session.Session
	greengrass    *greengrass.Greengrass
	iam           *iam.IAM
	iot           *iot.IoT
	lambda        *lambda.Lambda
	config        ggGroupConfig
	ggconfig      GreengrassConfig
	keyCertOutput *iot.CreateKeysAndCertificateOutput
}

// NewGreengrassSession is the constructor for Greengrass interactions
func NewGreengrassSession(sess *session.Session) *GreengrassSession {
	ggSession := new(GreengrassSession)
	ggSession.session = sess
	ggSession.greengrass = greengrass.New(sess)
	ggSession.iot = iot.New(sess)
	ggSession.iam = iam.New(sess)
	ggSession.lambda = lambda.New(sess)
	return ggSession
}

// WriteGGConfig - write the Greengrass Core config file
func (ggSession *GreengrassSession) WriteGGConfig(w io.Writer) {
	ggSession.ggconfig.Write(w)
}

// LoadGroupConfig - load the Greengrass group config file
func (ggSession *GreengrassSession) LoadGroupConfig(filename string) error {
	return ggSession.config.LoadConfig(filename)
}

// WriteGroupConfig - write the Greengrass group config file
func (ggSession *GreengrassSession) WriteGroupConfig(filename string) error {
	return ggSession.config.WriteConfig(filename)
}

// Cleanup - delete configuration for a Greengrass session object
func (ggSession *GreengrassSession) Cleanup() error {
	inactive := "INACTIVE"
	ggSession.iot.UpdateCertificate(&iot.UpdateCertificateInput{
		CertificateId: ggSession.keyCertOutput.CertificateId,
		NewStatus:     &inactive,
		//NewStatus: &iot.CertificateStatusInactive,
	})

	_, err := ggSession.iot.DetachThingPrincipal(&iot.DetachThingPrincipalInput{
		Principal: ggSession.keyCertOutput.CertificateId,
		ThingName: &ggSession.config.Core.ThingName,
	})
	if err != nil {
		fmt.Printf("DetachPrincipalPolicy error: %v\n", err)
		return err
	}

	_, err = ggSession.iot.DeleteCertificate(&iot.DeleteCertificateInput{
		CertificateId: ggSession.keyCertOutput.CertificateId,
	})
	if err != nil {
		fmt.Printf("cleanup: %+v\n", err)
		return err
	}

	return nil
}
