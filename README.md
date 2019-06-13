# Configor

Golang Configuration tool that support YAML, JSON, TOML, Shell Environment (Supports Go 1.10+)

[![wercker status](https://app.wercker.com/status/9ebd3684ff8998501af5aac38a79a380/s/master "wercker status")](https://app.wercker.com/project/byKey/9ebd3684ff8998501af5aac38a79a380)

## Usage

```go
package main

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var Config = struct {
	APPName string `default:"app name"`

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
}{}

func main() {
	configor.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}
```

With configuration file *config.yml*:

```yaml
appname: test

db:
    name:     test
    user:     test
    password: test
    port:     1234

contacts:
- name: i test
  email: test@test.com
```

## Debug Mode & Verbose Mode

Debug/Verbose mode is helpful when debuging your application, `debug mode` will let you know how `configor` loaded your configurations, like from which file, shell env, `verbose mode` will tell you even more, like those shell environments `configor` tried to load.

```go
// Enable debug mode or set env `CONFIGOR_DEBUG_MODE` to true when running your application
configor.New(&configor.Config{Debug: true}).Load(&Config, "config.json")

// Enable verbose mode or set env `CONFIGOR_VERBOSE_MODE` to true when running your application
configor.New(&configor.Config{Verbose: true}).Load(&Config, "config.json")
```

## Auto Reload Mode

Configor can auto reload configuration based on time

```go
// auto reload configuration every second
configor.New(&configor.Config{AutoReload: true}).Load(&Config, "config.json")

// auto reload configuration every minute
configor.New(&configor.Config{AutoReload: true, AutoReloadInterval: time.Minute}).Load(&Config, "config.json")
```

Auto Reload Callback

```go
configor.New(&configor.Config{AutoReload: true, AutoReloadCallback: func(config interface{}) {
    fmt.Printf("%v changed", config)
}}).Load(&Config, "config.json")
```

# Advanced Usage

* Load mutiple configurations

```go
// Earlier configurations have higher priority
configor.Load(&Config, "application.yml", "database.json")
```

* Return error on unmatched keys

Return an error on finding keys in the config file that do not match any fields in the config struct.
In the example below, an error will be returned if config.toml contains keys that do not match any fields in the ConfigStruct struct.
If ErrorOnUnmatchedKeys is not set, it defaults to false.

Note that for json files, setting ErrorOnUnmatchedKeys to true will have an effect only if using go 1.10 or later.

```go
err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&ConfigStruct, "config.toml")
```

* Load configuration by environment

Use `CONFIGOR_ENV` to set environment, if `CONFIGOR_ENV` not set, environment will be `development` by default, and it will be `test` when running tests with `go test`

```go
// config.go
configor.Load(&Config, "config.json")

$ go run config.go
// Will load `config.json`, `config.development.json` if it exists
// `config.development.json` will overwrite `config.json`'s configuration
// You could use this to share same configuration across different environments

$ CONFIGOR_ENV=production go run config.go
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration

$ go test
// Will load `config.json`, `config.test.json` if it exists
// `config.test.json` will overwrite `config.json`'s configuration

$ CONFIGOR_ENV=production go test
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration
```

```go
// Set environment by config
configor.New(&configor.Config{Environment: "production"}).Load(&Config, "config.json")
```

* Example Configuration

```go
// config.go
configor.Load(&Config, "config.yml")

$ go run config.go
// Will load `config.example.yml` automatically if `config.yml` not found and print warning message
```

* Load From Shell Environment

```go
$ CONFIGOR_APPNAME="hello world" CONFIGOR_DB_NAME="hello world" go run config.go
// Load configuration from shell environment, it's name is {{prefix}}_FieldName
```

```go
// You could overwrite the prefix with environment CONFIGOR_ENV_PREFIX, for example:
$ CONFIGOR_ENV_PREFIX="WEB" WEB_APPNAME="hello world" WEB_DB_NAME="hello world" go run config.go

// Set prefix by config
configor.New(&configor.Config{ENVPrefix: "WEB"}).Load(&Config, "config.json")
```

* Anonymous Struct

Add the `anonymous:"true"` tag to an anonymous, embedded struct to NOT include the struct name in the environment
variable of any contained fields.  For example:

```go
type Details struct {
	Description string
}

type Config struct {
	Details `anonymous:"true"`
}
```

With the `anonymous:"true"` tag specified, the environment variable for the `Description` field is `CONFIGOR_DESCRIPTION`.
Without the `anonymous:"true"`tag specified, then environment variable would include the embedded struct name and be `CONFIGOR_DETAILS_DESCRIPTION`.

* With flags

```go
func main() {
	config := flag.String("file", "config.yml", "configuration file")
	flag.StringVar(&Config.APPName, "name", "", "app name")
	flag.StringVar(&Config.DB.Name, "db-name", "", "database name")
	flag.StringVar(&Config.DB.User, "db-user", "root", "database user")
	flag.Parse()

	os.Setenv("CONFIGOR_ENV_PREFIX", "-")
	configor.Load(&Config, *config)
	// configor.Load(&Config) // only load configurations from shell env & flag
}
```

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

## Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the MIT License
