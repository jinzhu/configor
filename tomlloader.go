package configor

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Tomloader used to load / dump TOML files
type Tomlloader struct{}

// Load will read the file and unmarshal
func (l *Tomlloader) Load(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(file, ".toml"):
		return toml.Unmarshal(data, config)
	}
	return nil
}

// Dump will marshal config to a file
func (l *Tomlloader) Dump(config interface{}, file string) error {
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(config); err == nil {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()
		f.Write(buffer.Bytes())
	} else {
		return err
	}
	return nil
}
