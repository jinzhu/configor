package loader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

// Tomloader used to load / dump TOML files
type Tomlloader struct{}

// Load will read the file and unmarshal
func (l *Tomlloader) Load(config interface{}, file string) error {
	if !strings.HasSuffix(file, ".toml") {
		return fmt.Errorf("File does not have the toml extension: %s", file)
	}
	return l.PlainLoad(config, file)
}

// PlainLoad just does the unmarshalling
func (l *Tomlloader) PlainLoad(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, config)
}

// Dump will marshal config to a file
func (l *Tomlloader) Dump(config interface{}, file string) error {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
