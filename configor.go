package configor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

func ENV() string {
	if env := os.Getenv("CONFIGOR_ENV"); env != "" {
		return env
	}
	if isTest, _ := regexp.MatchString("/_test/", os.Args[0]); isTest {
		return "test"
	}
	return "development"
}

func getConfigurationWithENV(file, env string) (string, error) {
	var envFile string
	var extname = path.Ext(file)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, nil
	}
	return "", fmt.Errorf("failed to find file %v", file)
}

func getConfigurations(files ...string) []string {
	var results []string
	env := ENV()
	for _, file := range files {
		var foundFile bool
		// check configuration
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			results = append(results, file)
		}

		// check env configuration
		if file, err := getConfigurationWithENV(file, env); err == nil {
			foundFile = true
			results = append(results, file)
		}

		// check example configuration
		if !foundFile {
			if example, err := getConfigurationWithENV(file, "example"); err == nil {
				fmt.Printf("Failed to find configuration %v, using example file %v\n", file, example)
				results = append(results, file)
			} else {
				fmt.Printf("Failed to find configuration %v\n", file)
			}
		}
	}
	return results
}

func Load(config interface{}, files ...string) error {
	for _, file := range getConfigurations(files...) {
		if err := load(config, file); err != nil {
			return err
		}
	}

	return processTags(config)
}

func processTags(config interface{}) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		fieldStruct := configType.Field(i)
		field := configValue.Field(i)

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
			if value := fieldStruct.Tag.Get("default"); value != "" {
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
			} else if fieldStruct.Tag.Get("required") == "true" {
				return errors.New(fieldStruct.Name + " is required, but blank")
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := processTags(field.Addr().Interface()); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			var length = field.Len()
			for i := 0; i < length; i++ {
				if err := processTags(field.Index(i).Addr().Interface()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func load(config interface{}, file string) error {
	if data, err := ioutil.ReadFile(file); err == nil {
		switch {
		case strings.HasSuffix(file, ".yaml"), strings.HasSuffix(file, ".yml"):
			return yaml.Unmarshal(data, config)
		case strings.HasSuffix(file, ".json"):
			return json.Unmarshal(data, config)
		case strings.HasSuffix(file, ".ini"):
			return ini.MapTo(config, file)
		default:
			if json.Unmarshal(data, config) != nil {
				if yaml.Unmarshal(data, config) != nil {
					if ini.MapTo(config, file) != nil {
						return errors.New("failed to load file")
					}
				}
			}
			return nil
		}
	} else {
		return err
	}
}
