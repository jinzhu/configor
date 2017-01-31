package loader

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v1"
)

// Yamlloader used to load / dump YAML files
type Yamlloader struct{}

// Load will read the file and unmarshal
func (l *Yamlloader) Load(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		return yaml.Unmarshal(data, config)
	}
	return nil
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
