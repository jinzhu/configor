# Configor

Golang Configuration tool that support YAML, JSON, Shell Environment

# Usage

```go
import (
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

configor.Load(&Config, "config.yml", "config.json"...)
```

# Advanced Usage

* Different configuration for each environment

Use `CONFIGOR_ENV` to set the environment

```go
// config.go
configor.Load(&Config, "config.json")

$ CONFIGOR_ENV=production go run config.go
// Will load `config.yml`, `config.production.yml` if it is exist
// And `config.production.yml` will overwrite `config.yml`'s configuration
// You could use this to share same configuration across different environments
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

# Author

**jinzhu**

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

## License

Released under the MIT License
