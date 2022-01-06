package cfg

import (
	"context"
	"strconv"

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
	f, err := cfg.File(cfg.ParseYAML(), "config.yaml")
	if err != nil {
		return nil, err
	}

	c := &Config{
		Server: &svr.Config{
			Port: cfg.Schema[int]{
				Name:      "server port",
				Default:   func() int { return 8080 },
				Validator: cfg.Between(5000, 9000),
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
					cfg.Convert(cfg.Float64ToInt, f.Float64("server", "port")),
				},
			},
			EnableSSL: cfg.Schema[bool]{
				Name: "server enable ssl",
				Providers: []cfg.Provider[bool]{
					cfg.EnvVar("APP_ENABLE_SSL", strconv.ParseBool),
					f.Bool("server", "enable_ssl"),
				},
			},
		},
		DB: &db.Config{
			Host: cfg.Schema[string]{
				Name:    "db host",
				Default: func() string { return "localhost" },
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_HOST"),
					f.String("db", "host"),
				},
			},
			Port: cfg.Schema[int]{
				Name:    "db port",
				Default: func() int { return 3306 },
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
					cfg.Convert(cfg.Float64ToInt, f.Float64("db", "port")),
				},
			},
			Username: cfg.Schema[string]{
				Name:      "db username",
				Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_USERNAME"),
					f.String("db", "username"),
				},
			},
			Password: cfg.Schema[string]{
				Name: "db password",
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_PASSWORD"),
					f.String("db", "password"),
				},
			},
		},
	}

	if err := c.Validate(ctx); err != nil {
		return nil, err
	}

	return c, nil
}
