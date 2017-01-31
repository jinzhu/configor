package loader_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gophersgang/configor/loader"
)

var (
	yamlConfigFile = "/tmp/yaml_config.yaml"
)

func TestYamlDump(t *testing.T) {
	config := generateDefaultConfig()
	loader := &loader.Yamlloader{}
	loader.Dump(config, yamlConfigFile)

	dat, err := ioutil.ReadFile(yamlConfigFile)
	if err != nil {
		t.Error(err)
	}
	_ = dat
}

func TestYamlLoad(t *testing.T) {
	config := generateDefaultConfig()
	little := Config{
		APPName: "little config",
	}
	loader := &loader.Yamlloader{}

	loader.Dump(config, yamlConfigFile)
	loader.Load(&little, yamlConfigFile)

	if little.APPName != "configor" {
		t.Errorf("expected AppName to be configor, was %s", little.APPName)
	}
	fmt.Println(little)
}
