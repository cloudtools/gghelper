package gghelper_test

import (
	"bytes"
	"testing"

	"github.com/cloudtools/gghelper"
)

var expected = `{
  "coreThing": {
    "caPath": "foobar",
    "certPath": "",
    "keyPath": "",
    "thingArn": "",
    "iotHost": "",
    "ggHost": ""
  },
  "runtime": {
    "cgroup": {
      "useSystemd": ""
    }
  },
  "managedRespawn": false
}`

func TestGGConfig(t *testing.T) {
	ggc := gghelper.NewGGConfig()
	ggc.CoreThing.CAPath = "foobar"

	buf := new(bytes.Buffer)
	ggc.Write(buf)
	if buf.String() != expected {
		t.Errorf(" got: '%s'\nwant: '%s'\n", buf.String(), expected)
	}
}
