package configor_test

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/jinzhu/configor"
)

type Config struct {
	APPName string `default:"configor"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}

func generateDefaultConfig() Config {
	config := Config{
		APPName: "configor",
		DB: struct {
			Name     string
			User     string `default:"root"`
			Password string `required:"true"`
			Port     uint   `default:"3306"`
		}{
			Name:     "configor",
			User:     "configor",
			Password: "configor",
			Port:     3306,
		},
		Contacts: []struct {
			Name  string
			Email string `required:"true"`
		}{
			{
				Name:  "Jinzhu",
				Email: "wosmvp@gmail.com",
			},
		},
	}
	return config
}

func TestLoadNormalConfig(t *testing.T) {
	config := generateDefaultConfig()
	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			file.WriteString(string(bytes))
			var result Config
			configor.Load(&result, file.Name())
			if !reflect.DeepEqual(result, config) {
				t.Errorf("result should equal to original configuration")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestDefaultValue(t *testing.T) {
	config := generateDefaultConfig()
	config.APPName = ""
	config.DB.Port = 0

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			file.WriteString(string(bytes))
			var result Config
			configor.Load(&result, file.Name())
			if !reflect.DeepEqual(result, generateDefaultConfig()) {
				t.Errorf("result should be set default value correctly")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestMissingRequiredValue(t *testing.T) {
	config := generateDefaultConfig()
	config.DB.Password = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			file.WriteString(string(bytes))
			var result Config
			if err := configor.Load(&result, file.Name()); err == nil {
				t.Errorf("Should got error when load configuration missing db password")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}
