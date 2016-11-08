package configor_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/jinzhu/configor"
	"github.com/BurntSushi/toml"
)

type Config struct {
	APPName string `default:"configor"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
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
			Password string `required:"true" env:"DBPassword"`
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
			defer os.Remove(file.Name())
			file.Write(bytes)
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

func TestLoadTOMLConfigWithTomlExtension(t *testing.T) {
	config := generateDefaultConfig()
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor.toml"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(buffer.Bytes())
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

func TestLoadTOMLConfigWithoutExtension(t *testing.T) {
	config := generateDefaultConfig()
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(buffer.Bytes())
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
			defer os.Remove(file.Name())
			file.Write(bytes)
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
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			if err := configor.Load(&result, file.Name()); err == nil {
				t.Errorf("Should got error when load configuration missing db password")
			}
		}
	} else {
		t.Errorf("failed to marshal config")
	}
}

func TestLoadConfigurationByEnvironment(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "config2",
	}

	if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
		defer file.Close()
		defer os.Remove(file.Name())
		configBytes, _ := yaml.Marshal(config)
		config2Bytes, _ := yaml.Marshal(config2)
		ioutil.WriteFile(file.Name()+".yaml", configBytes, 0644)
		defer os.Remove(file.Name() + ".yaml")
		ioutil.WriteFile(file.Name()+".production.yaml", config2Bytes, 0644)
		defer os.Remove(file.Name() + ".production.yaml")

		var result Config
		os.Setenv("CONFIGOR_ENV", "production")
		defer os.Setenv("CONFIGOR_ENV", "")
		if err := configor.Load(&result, file.Name()+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		var defaultConfig = generateDefaultConfig()
		defaultConfig.APPName = "config2"
		if !reflect.DeepEqual(result, defaultConfig) {
			t.Errorf("result should be load configurations by environment correctly")
		}
	}
}

func TestOverwriteConfigurationWithEnvironmentWithDefaultPrefix(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			os.Setenv("CONFIGOR_APPNAME", "config2")
			os.Setenv("CONFIGOR_DB_NAME", "db_name")
			defer os.Setenv("CONFIGOR_APPNAME", "")
			defer os.Setenv("CONFIGOR_DB_NAME", "")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestOverwriteConfigurationWithEnvironment(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			os.Setenv("CONFIGOR_ENV_PREFIX", "app")
			os.Setenv("APP_APPNAME", "config2")
			os.Setenv("APP_DB_NAME", "db_name")
			defer os.Setenv("CONFIGOR_ENV_PREFIX", "")
			defer os.Setenv("APP_APPNAME", "")
			defer os.Setenv("APP_DB_NAME", "")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestResetPrefixToBlank(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			os.Setenv("CONFIGOR_ENV_PREFIX", "-")
			os.Setenv("APPNAME", "config2")
			os.Setenv("DB_NAME", "db_name")
			defer os.Setenv("CONFIGOR_ENV_PREFIX", "")
			defer os.Setenv("APPNAME", "")
			defer os.Setenv("DB_NAME", "")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.DB.Name = "db_name"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestReadFromEnvironmentWithSpecifiedEnvName(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			os.Setenv("DBPassword", "db_password")
			defer os.Setenv("DBPassword", "db_password")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.DB.Password = "db_password"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestENV(t *testing.T) {
	if configor.ENV() != "test" {
		t.Errorf("Env should be test when running `go test`")
	}

	os.Setenv("CONFIGOR_ENV", "production")
	defer os.Setenv("CONFIGOR_ENV", "")
	if configor.ENV() != "production" {
		t.Errorf("Env should be production when set it with CONFIGOR_ENV")
	}
}

func TestErrorIncludesPath(t *testing.T) {
	// Used in the array case
	type arrayStruct struct {
		B string `required:"true"`
	}

	cases := []struct {
		path   string
		strukt interface{}
	}{
		{
			// Top-level field
			path: "A",
			strukt: &struct {
				A string `required:"true"`
			}{},
		},
		{
			// Nested struct
			path: "A.B",
			strukt: &struct {
				A struct {
					B string `required:"true"`
				}
			}{},
		},
		{
			// Slice
			path: "A.0.B",
			strukt: &struct {
				A []arrayStruct
			}{
				A: []arrayStruct{
					{B: ""},
				},
			},
		},
	}

	for _, c := range cases {
		err := configor.Load(c.strukt)

		if err == nil {
			t.Error("no error on missing required field")
		}

		if !strings.Contains(err.Error(), c.path) {
			t.Errorf("error message doesn't contain path %s: %v", c.path, err)
		}
	}
}
