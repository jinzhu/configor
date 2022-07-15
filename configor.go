package configor

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Configor struct {
	*Config
	configModTimes map[string]time.Time
}

type Config struct {
	Environment        string
	ENVPrefix          string
	Debug              bool
	Verbose            bool
	Silent             bool
	AutoReload         bool
	AutoReloadInterval time.Duration
	AutoReloadCallback func(config interface{})

	// In case of json files, this field will be used only when compiled with
	// go 1.10 or later.
	// This field will be ignored when compiled with go versions lower than 1.10.
	ErrorOnUnmatchedKeys bool

	OnConfigChange func(e fsnotify.Event)
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

	if os.Getenv("CONFIGOR_SILENT_MODE") != "" {
		config.Silent = true
	}

	if config.AutoReload && config.AutoReloadInterval == 0 {
		config.AutoReloadInterval = time.Second
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
func (configor *Configor) Load(config interface{}, files ...string) (err error) {
	defaultValue := reflect.Indirect(reflect.ValueOf(config))
	if !defaultValue.CanAddr() {
		return fmt.Errorf("Config %v should be addressable", config)
	}
	err, _ = configor.load(config, false, files...)

	if configor.Config.AutoReload {
		go func() {
			timer := time.NewTimer(configor.Config.AutoReloadInterval)
			for range timer.C {
				reflectPtr := reflect.New(reflect.ValueOf(config).Elem().Type())
				reflectPtr.Elem().Set(defaultValue)

				var changed bool
				if err, changed = configor.load(reflectPtr.Interface(), true, files...); err == nil && changed {
					reflect.ValueOf(config).Elem().Set(reflectPtr.Elem())
					if configor.Config.AutoReloadCallback != nil {
						configor.Config.AutoReloadCallback(config)
					}
				} else if err != nil {
					fmt.Printf("Failed to reload configuration from %v, got error %v\n", files, err)
				}
				timer.Reset(configor.Config.AutoReloadInterval)
			}
		}()
	}
	return
}

// ENV return environment
func ENV() string {
	return New(nil).GetEnvironment()
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}

func Watch(config interface{}, files ...string) {
	New(nil).Watch(config, files...)
}

func (configor *Configor) Watch(config interface{}, files ...string) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	for _, file := range files {

		go func() {
			watcher, err := newWatcher()
			if err != nil {
				fmt.Printf("Failed to create watcher for %v, got error %v\n", file, err)
			}
			defer watcher.Close()

			configFile := filepath.Clean(file)
			configDir, _ := filepath.Split(configFile)
			realConfigFile, _ := filepath.EvalSymlinks(file)

			eventsWG := sync.WaitGroup{}
			eventsWG.Add(1)
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							eventsWG.Done()
							return
						}

						currentConfigFile, _ := filepath.EvalSymlinks(file)
						const writeOrCreateMask = fsnotify.Write | fsnotify.Create
						if (event.Op&writeOrCreateMask != 0 && currentConfigFile == event.Name) ||
							(currentConfigFile != "" && currentConfigFile != realConfigFile) {
							err := configor.Load(config, files...)
							if err != nil {
								fmt.Printf("Failed to reload configuration from %v, got error %v\n", files, err)
							}
							if configor.Config.OnConfigChange != nil {
								configor.Config.OnConfigChange(event)
							}
						} else if filepath.Clean(event.Name) == configFile &&
							event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
							eventsWG.Done()
							return
						}

					case err, ok := <-watcher.Errors:
						if ok {
							fmt.Printf("Failed to watch configuration from %v, got error %v\n", files, err)
						}
						eventsWG.Done()
						return
					}
				}
			}()
			watcher.Add(configDir)
			initWG.Done()
			eventsWG.Wait()
		}()
		initWG.Wait()
	}
}
