package configor_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/minitauros/configor"
	"gopkg.in/yaml.v2"
)

type Anonymous struct {
	Description string
}

type Config struct {
	APPName string `default:"configor"`
	Hosts   []string

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string `default:"Jinzhu"`
		Email string `required:"true"`
	}

	Anonymous `anonymous:"true"`

	MyMap map[string]string

	private string
}

func generateDefaultConfig() Config {
	config := Config{
		APPName: "configor",
		Hosts:   []string{"http://example.org", "http://jinzhu.me"},
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
			Name  string `default:"Jinzhu"`
			Email string `required:"true"`
		}{
			{
				Name:  "Jinzhu",
				Email: "wosmvp@gmail.com",
			},
		},
		Anonymous: Anonymous{
			Description: "This is an anonymous embedded struct whose environment variables should NOT include 'ANONYMOUS'",
		},
		MyMap: map[string]string{"a": "b"},
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

func TestLoadConfigFromTomlWithExtension(t *testing.T) {
	var (
		config = generateDefaultConfig()
		buffer bytes.Buffer
	)

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

func TestLoadConfigFromTomlWithoutExtension(t *testing.T) {
	var (
		config = generateDefaultConfig()
		buffer bytes.Buffer
	)

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
	config.Contacts[0].Name = ""

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())

			byteString := string(bytes)

			// Remove the appname setting from the json, for if we specify it (even a zero-value) it will overwrite the default.
			r := regexp.MustCompile(`"APPName":"",`)
			byteString = r.ReplaceAllString(byteString, "")

			file.Write([]byte(byteString))

			var result Config
			configor.Load(&result, file.Name())

			// The value for APPName should have been set, as its value was removed completely from the loaded JSON.
			// DB port however should not have been overwritten by the default value, as we specified a value (a zero-value, but we did specify somoething) for it in our config.
			defConfig := generateDefaultConfig()
			defConfig.DB.Port = 0

			if !reflect.DeepEqual(result, defConfig) {
				t.Errorf("default values were not correctly set")
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

func TestUnmatchedKeyInTomlConfigFile(t *testing.T) {
	type configStruct struct {
		Name string
	}
	type configFile struct {
		Name string
		Test string
	}
	config := configFile{Name: "test", Test: "ATest"}

	file, err := ioutil.TempFile("/tmp", "configor")
	if err != nil {
		t.Fatal("Could not create temp file")
	}
	defer os.Remove(file.Name())
	defer file.Close()

	filename := file.Name()

	if err := toml.NewEncoder(file).Encode(config); err == nil {

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := configor.New(&configor.Config{}).Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys")
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename)
		if err == nil {
			t.Errorf("Should get error when loading configuration with extra keys")
		}

		// The error should be of type UnmatchedTomlKeysError
		tomlErr, ok := err.(*configor.UnmatchedTomlKeysError)
		if !ok {
			t.Errorf("Should get UnmatchedTomlKeysError error when loading configuration with extra keys")
		}

		// The error.Keys() function should return the "Test" key
		keys := configor.GetStringTomlKeys(tomlErr.Keys)
		if len(keys) != 1 || keys[0] != "Test" {
			t.Errorf("The UnmatchedTomlKeysError should contain the Test key")
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .toml to the file name and test again
	err = os.Rename(filename, filename+".toml")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename = filename + ".toml"
	defer os.Remove(filename)

	var result configStruct

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := configor.New(&configor.Config{}).Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	err = configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename)
	if err == nil {
		t.Errorf("Should get error when loading configuration with extra keys")
	}

	// The error should be of type UnmatchedTomlKeysError
	tomlErr, ok := err.(*configor.UnmatchedTomlKeysError)
	if !ok {
		t.Errorf("Should get UnmatchedTomlKeysError error when loading configuration with extra keys")
	}

	// The error.Keys() function should return the "Test" key
	keys := configor.GetStringTomlKeys(tomlErr.Keys)
	if len(keys) != 1 || keys[0] != "Test" {
		t.Errorf("The UnmatchedTomlKeysError should contain the Test key")
	}

}

func TestUnmatchedKeyInYamlConfigFile(t *testing.T) {
	type configStruct struct {
		Name string
	}
	type configFile struct {
		Name string
		Test string
	}
	config := configFile{Name: "test", Test: "ATest"}

	file, err := ioutil.TempFile("/tmp", "configor")
	if err != nil {
		t.Fatal("Could not create temp file")
	}

	defer os.Remove(file.Name())
	defer file.Close()

	filename := file.Name()

	if data, err := yaml.Marshal(config); err == nil {
		file.WriteString(string(data))

		var result configStruct

		// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
		if err := configor.New(&configor.Config{}).Load(&result, filename); err != nil {
			t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
		}

		// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
		if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename); err == nil {
			t.Errorf("Should get error when loading configuration with extra keys")

			// The error should be of type *yaml.TypeError
		} else if _, ok := err.(*yaml.TypeError); !ok {
			// || !strings.Contains(err.Error(), "not found in struct") {
			t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
		}

	} else {
		t.Errorf("failed to marshal config")
	}

	// Add .yaml to the file name and test again
	err = os.Rename(filename, filename+".yaml")
	if err != nil {
		t.Errorf("Could not add suffix to file")
	}
	filename = filename + ".yaml"
	defer os.Remove(filename)

	var result configStruct

	// Do not return error when there are unmatched keys but ErrorOnUnmatchedKeys is false
	if err := configor.New(&configor.Config{}).Load(&result, filename); err != nil {
		t.Errorf("Should NOT get error when loading configuration with extra keys. Error: %v", err)
	}

	// Return an error when there are unmatched keys and ErrorOnUnmatchedKeys is true
	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&result, filename); err == nil {
		t.Errorf("Should get error when loading configuration with extra keys")

		// The error should be of type *yaml.TypeError
	} else if _, ok := err.(*yaml.TypeError); !ok {
		// || !strings.Contains(err.Error(), "not found in struct") {
		t.Errorf("Error should be of type yaml.TypeError. Instead error is %v", err)
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

func TestLoadConfigurationByEnvironmentSetByConfig(t *testing.T) {
	config := generateDefaultConfig()
	config2 := struct {
		APPName string
	}{
		APPName: "production_config2",
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
		var Configor = configor.New(&configor.Config{Environment: "production"})
		if Configor.Load(&result, file.Name()+".yaml"); err != nil {
			t.Errorf("No error should happen when load configurations, but got %v", err)
		}

		var defaultConfig = generateDefaultConfig()
		defaultConfig.APPName = "production_config2"
		if !reflect.DeepEqual(result, defaultConfig) {
			t.Errorf("result should be load configurations by environment correctly")
		}

		if Configor.GetEnvironment() != "production" {
			t.Errorf("configor's environment should be production")
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
			os.Setenv("CONFIGOR_HOSTS", "- http://example.org\n- http://jinzhu.me")
			os.Setenv("CONFIGOR_DB_NAME", "db_name")
			os.Setenv("CONFIGOR_MYMAP", `{"x":"y"}`)
			defer os.Setenv("CONFIGOR_APPNAME", "")
			defer os.Setenv("CONFIGOR_HOSTS", "")
			defer os.Setenv("CONFIGOR_DB_NAME", "")
			defer os.Setenv("CONFIGOR_MYMAP", "")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.APPName = "config2"
			defaultConfig.Hosts = []string{"http://example.org", "http://jinzhu.me"}
			defaultConfig.DB.Name = "db_name"
			defaultConfig.MyMap = map[string]string{"x":"y"}
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

func TestOverwriteConfigurationWithEnvironmentThatSetByConfig(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			os.Setenv("APP1_APPName", "config2")
			os.Setenv("APP1_DB_Name", "db_name")
			defer os.Setenv("APP1_APPName", "")
			defer os.Setenv("APP1_DB_Name", "")

			var result Config
			var Configor = configor.New(&configor.Config{ENVPrefix: "APP1"})
			Configor.Load(&result, file.Name())

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

func TestResetPrefixToBlank2(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			os.Setenv("CONFIGOR_ENV_PREFIX", "-")
			os.Setenv("APPName", "config2")
			os.Setenv("DB_Name", "db_name")
			defer os.Setenv("CONFIGOR_ENV_PREFIX", "")
			defer os.Setenv("APPName", "")
			defer os.Setenv("DB_Name", "")
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
			defer os.Setenv("DBPassword", "")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.DB.Password = "db_password"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestAnonymousStruct(t *testing.T) {
	config := generateDefaultConfig()

	if bytes, err := json.Marshal(config); err == nil {
		if file, err := ioutil.TempFile("/tmp", "configor"); err == nil {
			defer file.Close()
			defer os.Remove(file.Name())
			file.Write(bytes)
			var result Config
			os.Setenv("CONFIGOR_DESCRIPTION", "environment description")
			defer os.Setenv("CONFIGOR_DESCRIPTION", "")
			configor.Load(&result, file.Name())

			var defaultConfig = generateDefaultConfig()
			defaultConfig.Anonymous.Description = "environment description"
			if !reflect.DeepEqual(result, defaultConfig) {
				t.Errorf("result should equal to original configuration")
			}
		}
	}
}

func TestENV(t *testing.T) {
	if configor.ENV() != "test" {
		t.Errorf("Env should be test when running `go test`, instead env is %v", configor.ENV())
	}

	os.Setenv("CONFIGOR_ENV", "production")
	defer os.Setenv("CONFIGOR_ENV", "")
	if configor.ENV() != "production" {
		t.Errorf("Env should be production when set it with CONFIGOR_ENV")
	}
}
