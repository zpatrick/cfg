# CFG

[![Go Doc](https://godoc.org/github.com/zpatrick/cfg?status.svg)](https://godoc.org/github.com/zpatrick/cfg)
[![Build Status](https://github.com/zpatrick/cfg/actions/workflows/go.yaml/badge.svg?branch=main)](https://github.com/zpatrick/cfg/actions/workflows/go.yaml?query=branch%3Amain)

This package is designed to house a common set of configuration-related features & patterns for Golang services. This include:

- Support for multiple sources of configuration.
- Providing default values and validation logic for specific settings.
- Package API which encourages high cohesion and low coupling with the rest of the application.

## Usage

```go
package config

import (
	"context"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/envvar"
	"github.com/zpatrick/cfg/ini"
)

type Config struct {
	ServerPort       int
	ServerTimeout    time.Duration
	DatabaseAddress  string
	DatabaseUsername string
}

func LoadConfig(ctx context.Context, path string) (*Config, error) {
	iniFile, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := cfg.Load(ctx, cfg.Schemas{
		"database address": cfg.Schema[string]{
			Dest:    &c.DatabaseAddress,
			Default: cfg.Addr("localhost:3306"),
			Provider: envvar.New("APP_DATABASE_ADDR"),
		},
		"database username": cfg.Schema[string]{
			Dest:      &c.DatabaseUsername,
			Default:   cfg.Addr("readonly"),
			Validator: cfg.OneOf("admin", "readonly", "readwrite"),
			Provider: envvar.New("APP_DATABASE_USERNAME"),
		},
		"server port": cfg.Schema[int]{
			Dest:    &c.ServerPort,
			Default: cfg.Addr(8080),
			Validator: cfg.Between(8000, 9000),
			Provider: cfg.MultiProvider {
				envvar.Newf("APP_SERVER_PORT", strconv.Atoi),
				iniFile.Int("server", "port"),
			},
		},
		"server timeout": cfg.Schema[time.Duration]{
			Dest:      &c.ServerTimeout,
			Default:   cfg.Addr(time.Second * 30),
			Validator: cfg.Between(time.Second, time.Minute*5),
			Provider:  cfg.MultiProvider {
				envvar.Newf("APP_SERVER_TIMEOUT", time.ParseDuration),
				iniFile.Duration("server", "timeout"),
			},
		},
	}); err != nil {
		return nil, err
	}

	return c, nil
}
```


# Built in Providers

- [Environment Variables](https://pkg.go.dev/github.com/zpatrick/cfg#EnvVar)
- [Flags](https://pkg.go.dev/github.com/zpatrick/cfg#Flag)
- [INI Files](https://pkg.go.dev/github.com/zpatrick/cfg#INIFile)
- [TOML Files](https://pkg.go.dev/github.com/zpatrick/cfg#TOMLFile)

Please see the [Godoc](https://pkg.go.dev/github.com/zpatrick/cfg#example-YAML) example for YAML files.  

# Validation
A setting may specify a [Validator](https://pkg.go.dev/github.com/zpatrick/cfg#Validator) which will check whether or not a provided value is valid.
The built in validators are:

- [Between](https://pkg.go.dev/github.com/zpatrick/cfg#Between) - Ensures a given value is between the specified parameters.
- [OneOf](https://pkg.go.dev/github.com/zpatrick/cfg#OneOf) - Ensures a given value is one of the specified parameters.
- [Or](https://pkg.go.dev/github.com/zpatrick/cfg#Or) - Combines multiple Validators, ensures at least one passes.
- [And](https://pkg.go.dev/github.com/zpatrick/cfg#And) - Combines multiple Validators, ensures all pass.
- [Not](https://pkg.go.dev/github.com/zpatrick/cfg#Not) - Ensures a given validator does not pass.

# Advanced

## Custom Validation
A custom Validator must satisfy the [Validator](https://pkg.go.dev/github.com/zpatrick/cfg#Validator) interface.
The simplest way to achieve this is by using the [ValidatorFunc](https://pkg.go.dev/github.com/zpatrick/cfg#ValidatorFunc) type.

```go
cfg.Setting[string]{
	Default: cfg.Addr("name@email.com"),
	Validator: cfg.ValidatorFunc(func(addr string) error {
		_, err := mail.ParseAddr(addr)
		return err
	}),
}
```

## Custom Providers

```go
import (
	"net/mail"
)

// Create a helper function which wraps the underlying type.
func provideMailAddr(provider cfg.Provider[string]) cfg.Provider[*mail.Addr] {
	return cfg.ProviderFunc[*mail.Addr](func(ctx context.Context) (*mail.Addr, error) {
		// Get the underlying value from the given provider.
		raw, err := provider.Provide(ctx)
		if err != nil {
			return out, err
		}

		// Convert the underlying.
		return mail.ParseAddr(raw)
	})
}

// Use the helper function in your cfg.Setting definition.
func LoadConfig() (*Config, error) {
	yamlFile, err := cfg.YAMLFile("config.yaml")
	if err != nil {
		return nil, err
	}

	email := cfg.Setting[*mail.Addr]{
		Providers: []cfg.Provider[*mail.Addr]{
			provideMailAddr(yamlFile.String("email")),
		},
	}

	...
}
```
