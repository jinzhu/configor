package loader_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gophersgang/configor/loader"
)

var (
	jsonConfigFile = "/tmp/json_config.json"
)

func TestJsonDump(t *testing.T) {
	config := generateDefaultConfig()
	loader := &loader.Jsonloader{}
	loader.Dump(config, jsonConfigFile)

	dat, err := ioutil.ReadFile(jsonConfigFile)
	if err != nil {
		t.Error(err)
	}
	_ = dat
}

func TestJsonLoad(t *testing.T) {
	config := generateDefaultConfig()
	little := Config{
		APPName: "little config",
	}
	loader := &loader.Jsonloader{}

	loader.Dump(config, jsonConfigFile)
	loader.Load(&little, jsonConfigFile)

	if little.APPName != "configor" {
		t.Errorf("expected AppName to be configor, was %s", little.APPName)
	}
	fmt.Println(little)
}
