# How to use the config struct

## Load config file

from `helper/config/config.go`

We declare a new type `Config` struct in this file. This `Config` struct will hold all configuration variables of the application that we read from file or environment variables

## Declare new config variable

Example: We want to add the dapr configuration, open the file `config-dev.yml` and add the block in the below

```yaml
dapr:
  url: https://example.com
```

from `helper/config/config.go`

```go

type Config struct {
    ....
	DaprConfig DaprConfig `mapstructure:"dapr"`
}

type DaprConfig struct {
	Url string `mapstructure:"url"`
}
```

## Usage 

```go
package main

import (
	"fmt"
	"marketplace-svc/helper/config"
)

func main() {
	cfg := config.Init()

	fmt.Println(cfg.DaprConfig.Url)
}
```