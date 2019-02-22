# Toml Types
Adds support to marshal and unmarshal types not in the official TOML spec but very useful.

## Usage

```go
package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/configor/toml"
)

var Config = struct {
	Name string `toml:"name"`
	App struct {
		Server     string `toml:"server"`
		Timeout    toml.Duration `toml:"timeout"`
	}
}{}

func main() {
	os.Setenv("CONFIGOR_ENV_PREFIX", "TEST")
	os.Setenv("TEST_APP_SERVER", "just for test")
	configor.Load(&Config, "config.toml")
	fmt.Printf("config: %#v", Config)
}
```

With configuration file *config.toml*:

```toml
name = "toml"
[app]
server = "127.0.0.1:8080"
timeout = "10s"
```