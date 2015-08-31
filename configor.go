package configor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

func Load(config interface{}, files ...string) error {
	for _, file := range files {
		if err := load(config, file); err != nil {
			return err
		}
	}

	processTags(config)
	return nil
}

func processTags(config interface{}) error {
	configValue := reflect.ValueOf(config).Elem()
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		fieldStruct := configType.Field(i)
		field := configValue.Field(i)
		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
			if value := fieldStruct.Tag.Get("default"); value != "" {
				return json.Unmarshal([]byte(value), field.Interface())
			} else if fieldStruct.Tag.Get("required") == "true" {
				return errors.New(fieldStruct.Name + " is required, but blank")
			}
		}
	}
	return nil
}

func load(config interface{}, file string) error {
	if data, err := ioutil.ReadFile(file); err == nil {
		switch {
		case strings.HasSuffix(file, ".yaml"), strings.HasSuffix(file, ".yml"):
			return json.Unmarshal(data, config)
		case strings.HasSuffix(file, ".json"):
			return yaml.Unmarshal(data, config)
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
