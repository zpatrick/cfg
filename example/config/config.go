package config

import (
	"context"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
)

const DefaultConfigFilePath = "config.yaml"

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

func Load(ctx context.Context, configFilePath string) (*Config, error) {
	yamlFile, err := cfg.YAMLFile(configFilePath)
	if err != nil {
		return nil, err
	}

	settings := cfg.Settings{
		"server.port": cfg.Setting[int]{
			Default:   func() int { return 8080 },
			Validator: cfg.Between(5000, 9000),
			Providers: []cfg.Provider[int]{
				cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
				yamlFile.Int("server", "port"),
			},
		},
		"server.timeout": cfg.Setting[time.Duration]{
			Providers: []cfg.Provider[time.Duration]{
				cfg.EnvVar("APP_SERVER_TIMEOUT", time.ParseDuration),
				yamlFile.Duration("server", "timeout"),
			},
		},
		"db.host": cfg.Setting[string]{
			Default: func() string { return "localhost" },
			Providers: []cfg.Provider[string]{
				cfg.EnvVarStr("APP_DB_HOST"),
				yamlFile.String("database", "host"),
			},
		},
		"db.port": cfg.Setting[int]{
			Default: func() int { return 3306 },
			Providers: []cfg.Provider[int]{
				cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
				yamlFile.Int("database", "port"),
			},
		},
		"db.username": cfg.Setting[string]{
			Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
			Providers: []cfg.Provider[string]{
				cfg.EnvVarStr("APP_DB_USERNAME"),
				yamlFile.String("database", "username"),
			},
		},
		"db.password": cfg.Setting[string]{
			Providers: []cfg.Provider[string]{
				cfg.EnvVarStr("APP_DB_PASSWORD"),
				yamlFile.String("database", "password"),
			},
		},
	}

	if err := settings.Validate(ctx); err != nil {
		return nil, err
	}

	c := &Config{
		Server: ServerConfig{
			Port:    cfg.MustGet[int](ctx, settings["server.port"]),
			Timeout: cfg.MustGet[time.Duration](ctx, settings["server.timeout"]),
		},
		DB: DBConfig{
			Host:     cfg.MustGet[string](ctx, settings["db.host"]),
			Port:     cfg.MustGet[int](ctx, settings["db.port"]),
			Username: cfg.MustGet[string](ctx, settings["db.username"]),
			Password: cfg.MustGet[string](ctx, settings["db.password"]),
		},
	}

	return c, nil
}
