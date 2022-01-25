# CFG

[![Go Doc](https://godoc.org/github.com/zpatrick/cfg?status.svg)](https://godoc.org/github.com/zpatrick/cfg)

The cfg package is written to house common configuration features used by microservice applications.
This includes:

- Supporting multiple configuration sources - such as environment variables, configuration files, and command-line flags.
- Validation and defaults for specific configuration values.
- A type-safe package API using generics.

## Usage
_Click [here](https://github.com/zpatrick/cfg/tree/main/example) for a fully working example._

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
)

func main() {
	iniFile, err := cfg.INIFile("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	port := cfg.Setting[int]{
		Default: func() int { return 8080 },
		Providers: []cfg.Provider[int]{
			cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
			iniFile.Int("server", "port"),
		},
	}

	timeout := cfg.Setting[time.Duration]{
		Validator: cfg.Between(time.Second, time.Minute*2),
		Providers: []cfg.Provider[time.Duration]{
			cfg.EnvVar("APP_SERVER_TIMEOUT", time.ParseDuration),
			iniFile.Duration("server", "timeout"),
		},
	}

	ctx := context.Background()
	svr := http.Server{
		ReadTimeout:  timeout.MustGet(ctx),
		WriteTimeout: timeout.MustGet(ctx),
		IdleTimeout:  timeout.MustGet(ctx),
		Addr:         fmt.Sprintf(":%d", port.MustGet(ctx)),
		Handler:      http.NotFoundHandler(),
	}

	log.Println("listening on", svr.Addr)
	log.Fatal(svr.ListenAndServe())
}
```

TODO: Note how []Providers finds the first value found in the list

# Built in Providers

- [Environment Variables](https://pkg.go.dev/github.com/zpatrick/cfg#EnvVar)
- [Flags](https://pkg.go.dev/github.com/zpatrick/cfg#Flag)
- [INI Files](https://pkg.go.dev/github.com/zpatrick/cfg#INIFile)
- [JSON Files](https://pkg.go.dev/github.com/zpatrick/cfg#JSONFile)
- [YAML Files](https://pkg.go.dev/github.com/zpatrick/cfg#YAMLFile)


# Validation
A setting may specify a [Validator](https://pkg.go.dev/github.com/zpatrick/cfg#Validator) which will check whether or not a provided value is valid.
The built in validators are:

- [Between](https://pkg.go.dev/github.com/zpatrick/cfg#Between) - Ensures a given value is between the specified parameters.
- [OneOf](https://pkg.go.dev/github.com/zpatrick/cfg#OneOf) - Ensures a given value is one of the specified parameters.
- [Or](https://pkg.go.dev/github.com/zpatrick/cfg#Or) - Combines multiple Validators, ensures at least one passes.
- [And](https://pkg.go.dev/github.com/zpatrick/cfg#And) - Combines multiple Validators, ensures all pass.

# Advanced

## Custom Providers

```go
import (
	"net/mail"
)

// Create a helper function which wraps the underlying type.
func provideMailAddress(provider cfg.Provider[string]) cfg.Provider[*mail.Address] {
	return cfg.ProviderFunc[*mail.Address](func(ctx context.Context) (out *mail.Address, err error) {
		// Get the underlying value from the given provider.
		raw, err := provider.Provide(ctx)
		if err != nil {
			return out, err
		}

		// Convert the underlying.
		return mail.ParseAddress(raw)
	})
}

// Use the helper function in your cfg.Setting definition.
func LoadConfig() (*Config, error) {
	yamlFile, err := cfg.YAMLFile("config.yaml")
	if err != nil {
		return nil, err
	}

	email := cfg.Setting[*mail.Address]{
		Providers: []cfg.Provider[RoundedNumber]{
			provideMailAddress(yamlFile.String("email")),
		},
	}

	return nil
}
```

## Custom Validation
A custom Validator must satisfy the [Validator](https://pkg.go.dev/github.com/zpatrick/cfg#Validator) interface.
The simplest way to achieve this is by using the [ValidatorFunc](https://pkg.go.dev/github.com/zpatrick/cfg#ValidatorFunc) type.

```go
import (
  "net/mail"
)

var email = cfg.Setting[string]{
	Default: func() string { return "foo@bar.com" },
	Validator: cfg.ValidatorFunc(func(addr string) error {
		_, err := mail.ParseAddress(addr)
		return err
	}),
}
```
