package config

import (
	"context"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/example/db"
	"github.com/zpatrick/cfg/example/svr"
	"go.uber.org/multierr"
)

type Config struct {
	Server *svr.Config
	DB     *db.Config
}

func (c Config) Validate(ctx context.Context) error {
	return multierr.Combine(
		cfg.Validate(ctx, "server", c.Server),
		cfg.Validate(ctx, "db", c.DB),
	)
}

func Load(ctx context.Context) (*Config, error) {
	f, err := cfg.YAMLFile("config.yaml")
	if err != nil {
		return nil, err
	}

	// TODO: When a schema is defined in a struct but not instantiated (e.g. we forget to define server.Timeout here),
	// The validation step fails witha a NoValueProvided error but doesn't have the name so it's hard to tell which
	// config provider it's talking about.
	c := &Config{
		Server: &svr.Config{
			Port: cfg.Schema[int]{
				Name:      "server port",
				Default:   func() int { return 8080 },
				Validator: cfg.Between(5000, 9000),
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
					f.Int("server", "port"),
				},
			},
			EnableSSL: cfg.Schema[bool]{
				Name: "server enable ssl",
				Providers: []cfg.Provider[bool]{
					cfg.EnvVar("APP_ENABLE_SSL", strconv.ParseBool),
					f.Bool("server", "enable_ssl"),
				},
			},
			Timeout: cfg.Schema[time.Duration]{
				Name: "server timeout",
				Providers: []cfg.Provider[time.Duration]{
					cfg.EnvVar("APP_SERVER_TIMEOUT", time.ParseDuration),
					f.Duration("server", "timeout"),
				},
			},
		},
		DB: &db.Config{
			Host: cfg.Schema[string]{
				Name:    "db host",
				Default: func() string { return "localhost" },
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_HOST"),
					f.String("database", "host"),
				},
			},
			Port: cfg.Schema[int]{
				Name:    "db port",
				Default: func() int { return 3306 },
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
					f.Int("database", "port"),
				},
			},
			Username: cfg.Schema[string]{
				Name:      "db username",
				Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_USERNAME"),
					f.String("database", "username"),
				},
			},
			Password: cfg.Schema[string]{
				Name: "db password",
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_PASSWORD"),
					f.String("database", "password"),
				},
			},
		},
	}

	if err := c.Validate(ctx); err != nil {
		return nil, err
	}

	return c, nil
}
