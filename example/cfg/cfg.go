package cfg

import (
	"context"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/example/db"
	"github.com/zpatrick/cfg/example/svr"
)

type Config struct {
	Server *svr.Config
	DB     *db.Config
}

func (c Config) Validate(ctx context.Context) error {
	return cfg.Validate(ctx, c.Server, c.DB)
}

func Load(ctx context.Context) (*Config, error) {
	configFile, err := cfg.File(cfg.FormatYAML, "config.yaml")
	if err != nil {
		return nil, err
	}

	c := &Config{
		Server: &svr.Config{
			Port: cfg.Schema[int]{
				Name:      "server port",
				Decode:    cfg.DecodeInt,
				Default:   func() int { return 8080 },
				Validator: cfg.Between(5000, 9000),
				Providers: []cfg.Provider{
					configFile.Provide("server", "port"),
					cfg.EnvVar("APP_SERVER_PORT"),
				},
			},
			EnableSSL: cfg.Schema[bool]{
				Name:   "server enable ssl",
				Decode: cfg.DecodeBool,
				Providers: []cfg.Provider{
					configFile.Provide("server", "enable_ssl"),
					cfg.EnvVar("APP_ENABLE_SSL"),
				},
			},
		},
		DB: &db.Config{
			Host: cfg.Schema[string]{
				Name:    "db host",
				Decode:  cfg.DecodeString,
				Default: func() string { return "localhost" },
				Providers: []cfg.Provider{
					configFile.Provide("db", "host"),
					cfg.EnvVar("APP_DB_HOST"),
				},
			},
			Port: cfg.Schema[int]{
				Name:    "db port",
				Decode:  cfg.DecodeInt,
				Default: func() int { return 3306 },
				Providers: []cfg.Provider{
					configFile.Provide("db", "port"),
					cfg.EnvVar("APP_DB_PORT"),
				},
			},
			Username: cfg.Schema[string]{
				Name:      "db username",
				Decode:    cfg.DecodeString,
				Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
				Providers: []cfg.Provider{
					configFile.Provide("db", "username"),
					cfg.EnvVar("APP_DB_USERNAME"),
				},
			},
			Password: cfg.Schema[string]{
				Name:   "db password",
				Decode: cfg.DecodeString,
				Providers: []cfg.Provider{
					configFile.Provide("db", "password"),
					cfg.EnvVar("APP_DB_PASSWORD"),
				},
			},
		},
	}

	if err := c.Validate(ctx); err != nil {
		return nil, err
	}

	return c, nil
}
