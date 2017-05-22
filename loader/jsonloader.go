package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Yamlloader used to load / dump YAML files
type Jsonloader struct{}

// Load will read the file and unmarshal
func (l *Jsonloader) Load(config interface{}, file string) error {
	if !strings.HasSuffix(file, ".json") {
		return fmt.Errorf("File does not have the json extension: %s", file)
	}
	return l.PlainLoad(config, file)
}

// PlainLoad just does the unmarshalling
func (l *Jsonloader) PlainLoad(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// Dump will marshal config to a file
func (l *Jsonloader) Dump(config interface{}, file string) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, configBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
