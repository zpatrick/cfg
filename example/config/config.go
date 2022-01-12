package config

import (
	"context"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
	"go.uber.org/multierr"
)

type Schema struct {
	Server *ServerSchema
	DB     *DBSchema
}

func (c *Schema) Validate(ctx context.Context) error {
	return multierr.Combine(
		cfg.Validate(ctx, "server", c.Server),
		cfg.Validate(ctx, "db", c.DB),
	)
}

type ServerSchema struct {
	Port    cfg.Schema[int]
	Timeout cfg.Schema[time.Duration]
}

func (s ServerSchema) Validate(ctx context.Context) error {
	return multierr.Combine(
		cfg.Validate(ctx, "port", s.Port),
		cfg.Validate(ctx, "timeout", s.Timeout),
	)
}

type DBSchema struct {
	Host     cfg.Schema[string]
	Port     cfg.Schema[int]
	Username cfg.Schema[string]
	Password cfg.Schema[string]
}

func (d DBSchema) Validate(ctx context.Context) error {
	return multierr.Combine(
		cfg.Validate(ctx, "host", d.Host),
		cfg.Validate(ctx, "port", d.Port),
		cfg.Validate(ctx, "username", d.Username),
		cfg.Validate(ctx, "password", d.Password),
	)
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port    int
	Timeout time.Duration
}

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func Load(ctx context.Context) (Config, error) {
	f, err := cfg.YAMLFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	schema := &Schema{
		Server: &ServerSchema{
			Port: cfg.Schema[int]{
				Name:      "server port",
				Default:   func() int { return 8080 },
				Validator: cfg.Between(5000, 9000),
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
					f.Int("server", "port"),
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
		DB: &DBSchema{
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

	// Calling Validate allows us to safely use schema.MustLoad below.
	if err := schema.Validate(ctx); err != nil {
		return Config{}, err
	}

	c := Config{
		Server: ServerConfig{
			Port:    schema.Server.Port.MustLoad(ctx),
			Timeout: schema.Server.Timeout.MustLoad(ctx),
		},
		DB: DBConfig{
			Host:     schema.DB.Host.MustLoad(ctx),
			Port:     schema.DB.Port.MustLoad(ctx),
			Username: schema.DB.Username.MustLoad(ctx),
			Password: schema.DB.Password.MustLoad(ctx),
		},
	}

	return c, nil
}
