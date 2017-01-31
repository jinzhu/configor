package configor_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gophersgang/configor"
)

var (
	tomlConfigFile = "/tmp/tomlconfig.toml"
)

func TestTomlDump(t *testing.T) {
	config := generateDefaultConfig()
	loader := &configor.Tomlloader{}
	loader.Dump(config, tomlConfigFile)

	dat, err := ioutil.ReadFile(tomlConfigFile)
	if err != nil {
		t.Error(err)
	}
	_ = dat
}

func TestTomlLoad(t *testing.T) {
	config := generateDefaultConfig()
	little := Config{
		APPName: "little config",
	}
	loader := &configor.Tomlloader{}

	loader.Dump(config, tomlConfigFile)
	loader.Load(&little, tomlConfigFile)

	if little.APPName != "configor" {
		t.Errorf("expected AppName to be configor, was %s", little.APPName)
	}
	fmt.Println(little)
}
