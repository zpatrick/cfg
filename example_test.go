package cfg_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/envvar"
	"github.com/zpatrick/cfg/providers/yaml"
)

type Config struct {
	ServerPort       int
	ServerTimeout    time.Duration
	DatabaseAddress  string
	DatabaseUsername string
}

func LoadConfig(ctx context.Context, path string) (*Config, error) {
	yamlFile, err := yaml.New(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := cfg.Load(ctx, map[string]cfg.Loader{
		"server.port": cfg.Schema[int]{
			Dest:    &c.ServerPort,
			Default: cfg.Addr(8080),
			Provider: cfg.MultiProvider[int]{
				envvar.Newf("APP_SERVER_PORT", strconv.Atoi),
				yamlFile.Int("server", "port"),
			},
		},
		"server.timeout": cfg.Schema[time.Duration]{
			Dest:      &c.ServerTimeout,
			Default:   cfg.Addr(time.Second * 30),
			Validator: cfg.Between(time.Second, time.Minute*5),
			Provider: cfg.MultiProvider[time.Duration]{
				envvar.Newf("APP_SERVER_TIMEOUT", time.ParseDuration),
				yamlFile.Duration("server", "timeout"),
			},
		},
		"database.address": cfg.Schema[string]{
			Dest:    &c.DatabaseAddress,
			Default: cfg.Addr("localhost:3306"),
			Provider: cfg.MultiProvider[string]{
				envvar.New("APP_DATABASE_ADDR"),
				yamlFile.String("db", "address"),
			},
		},
		"database.username": cfg.Schema[string]{
			Dest:      &c.DatabaseUsername,
			Default:   cfg.Addr("readonly"),
			Validator: cfg.OneOf("readonly", "readwrite"),
			Provider: cfg.MultiProvider[string]{
				envvar.New("APP_DATABASE_USERNAME"),
				yamlFile.String("db", "username"),
			},
		},
	}); err != nil {
		return nil, err
	}

	return &c, nil
}

func Example() {
	ctx := context.Background()
	c, err := LoadConfig(ctx, "config.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println("config:", c)
}
