# Configor

Golang Configuration tool that support YAML, JSON, Shell Environment

# Usage

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

# Advanced Usage

* Load mutiple configurations

```go
// Earlier configurations have higher priority
configor.Load(&Config, "application.yml", "database.json")
```

* Different configuration for each environment

Use `CONFIGOR_ENV` to set the environment.

If `CONFIGOR_ENV` not set, when running tests with `go test`, the ENV will be `test`, otherwise, it will be `development`

```go
// config.go
configor.Load(&Config, "config.json")

$ go run config.go
// Will load `config.json`, `config.development.json` if it is exist
// And `config.development.json` will overwrite `config.json`'s configuration
// You could use this to share same configuration across different environments

$ CONFIGOR_ENV=production go run config.go
// Will load `config.json`, `config.production.json` if it is exist
// And `config.production.json` will overwrite `config.json`'s configuration

$ go test
// Will load `config.json`, `config.test.json` if it is exist
// And `config.test.json` will overwrite `config.json`'s configuration

$ CONFIGOR_ENV=production go test
// Will load `config.json`, `config.production.json` if it is exist
// And `config.production.json` will overwrite `config.json`'s configuration
```

* Example Configuration

```go
// config.go
configor.Load(&Config, "config.yml")

$ go run config.go
// Will load `config.example.yml` automatically if `config.yml` not found and print warning message
```

* Read From Shell Environment

```go
$ CONFIGOR_APPNAME="hello world" CONFIGOR_DB_NAME="hello world" go run config.go
// Will use shell environment's value if found with upcase of prefix (by default is CONFIGOR) + field name as key
// You could overwrite the prefix with environment CONFIGOR_ENV_PREFIX, for example:
$ CONFIGOR_ENV_PREFIX="WEB" WEB_APPNAME="hello world" WEB_DB_NAME="hello world" go run config.go
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

# Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the MIT License
