package loader

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v1"
)

// Yamlloader used to load / dump YAML files
type Yamlloader struct{}

// Load will read the file and unmarshal
func (l *Yamlloader) Load(config interface{}, file string) error {
	if !strings.HasSuffix(file, ".yml") && !strings.HasSuffix(file, ".yaml") {
		return fmt.Errorf("File does not have the yaml / yml extension: %s", file)
	}
	return l.PlainLoad(config, file)
}

// PlainLoad just does the unmarshalling
func (l *Yamlloader) PlainLoad(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, config)
}

// Dump will marshal config to a file
func (l *Yamlloader) Dump(config interface{}, file string) error {
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, configBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
