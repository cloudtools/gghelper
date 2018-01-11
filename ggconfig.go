package gghelper

import (
	"encoding/json"
	"io"
)

// GreengrassConfig holds the /greengrass/config/config.json structure
type GreengrassConfig struct {
	CoreThing struct {
		CAPath   string `json:"caPath"`
		CertPath string `json:"certPath"`
		KeyPath  string `json:"keyPath"`
		ThingArn string `json:"thingArn"`
		IOTHost  string `json:"iotHost"`
		GGHost   string `json:"ggHost"`
	} `json:"coreThing"`
	Runtime struct {
		Cgroup struct {
			UseSystemd string `json:"useSystemd"`
		} `json:"cgroup"`
	} `json:"runtime"`
	ManagedRespawn bool `json:"managedRespawn"`
}

// NewGGConfig - create a new Greengrass core config object
func NewGGConfig() (ggconfig *GreengrassConfig) {
	ggconfig = new(GreengrassConfig)
	return ggconfig
}

func (ggc *GreengrassConfig) Write(w io.Writer) {
	b, _ := json.MarshalIndent(*ggc, "", "  ")
	w.Write(b)
}
