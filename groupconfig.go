package gghelper

import (
	"encoding/json"
	"io/ioutil"
)

type definition struct {
	ID         string `json:"id"`
	VersionArn string `json:"version_arn"`
}

type function struct {
	Arn          string `json:"arn"`
	ArnQualifier string `json:"arn_qualifier"`
}

type group struct {
	Arn     string `json:"arn"`
	ID      string `json:"id"`
	Version string `json:"version"`
}

type thing struct {
	CertArn   string `json:"cert_arn"`
	CertID    string `json:"cert_id"`
	ThingArn  string `json:"thing_arn"`
	ThingName string `json:"thing_name"`
}

type ggGroupConfig struct {
	Core                   thing               `json:"core"`
	CoreDefinition         definition          `json:"core_def"`
	DeviceDefinition       definition          `json:"device_def"`
	Devices                map[string]thing    `json:"device,omitempty"`
	FunctionDefinition     definition          `json:"func_def"`
	Group                  group               `json:"group"`
	LambdaFunctions        map[string]function `json:"lambda_functions,omitempty"`
	LoggerDefinition       definition          `json:"logger_def"`
	SubscriptionDefinition definition          `json:"subscription_def"`
}

func (group *ggGroupConfig) LoadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	json.Unmarshal(data, group)
	if group.LambdaFunctions == nil {
		group.LambdaFunctions = make(map[string]function)
	}

	return nil
}

func (group *ggGroupConfig) WriteConfig(filename string) error {
	data, err := json.MarshalIndent(group, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (group *ggGroupConfig) ClearConfig() {
	newConfig := new(ggGroupConfig)
	group = newConfig
}
