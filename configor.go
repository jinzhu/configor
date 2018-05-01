package configor

import (
	"fmt"
	"os"
	"regexp"
)

type Configor struct {
	*Config
}

type Config struct {
	Environment string
	ENVPrefix   string
	Debug       bool
	Verbose     bool

	// In case of json files, this field will be used only when compiled with
	// go 1.10 or later.
	// This field will be ignored when compiled with go versions lower than 1.10.
	ErrorOnUnmatchedKeys bool
}

// New initialize a Configor
func New(config *Config) *Configor {
	if config == nil {
		config = &Config{}
	}

	if os.Getenv("CONFIGOR_DEBUG_MODE") != "" {
		config.Debug = true
	}

	if os.Getenv("CONFIGOR_VERBOSE_MODE") != "" {
		config.Verbose = true
	}

	return &Configor{Config: config}
}

var testRegexp = regexp.MustCompile("_test|(\\.test$)")

// GetEnvironment get environment
func (configor *Configor) GetEnvironment() string {
	if configor.Environment == "" {
		if env := os.Getenv("CONFIGOR_ENV"); env != "" {
			return env
		}

		if testRegexp.MatchString(os.Args[0]) {
			return "test"
		}

		return "development"
	}
	return configor.Environment
}

// GetErrorOnUnmatchedKeys returns a boolean indicating if an error should be
// thrown if there are keys in the config file that do not correspond to the
// config struct
func (configor *Configor) GetErrorOnUnmatchedKeys() bool {
	return configor.ErrorOnUnmatchedKeys
}

// Load will unmarshal configurations to struct from files that you provide
func (configor *Configor) Load(config interface{}, files ...string) error {
	defer func() {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Configuration:\n  %#v\n", config)
		}
	}()

	for _, file := range configor.getConfigurationFiles(files...) {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Loading configurations from file '%v'...\n", file)
		}
		if err := processFile(config, file, configor.GetErrorOnUnmatchedKeys()); err != nil {
			return err
		}
	}

	prefix := configor.getENVPrefix(config)
	if prefix == "-" {
		return configor.processTags(config)
	}
	return configor.processTags(config, prefix)
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}
