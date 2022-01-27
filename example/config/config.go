package config

import (
	"context"
	"strconv"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/example/database"
	"github.com/zpatrick/cfg/example/server"
)

const DefaultConfigFilePath = "config.yaml"

type Config struct {
	Server server.Config
	DB     database.Config
}

type ServerConfig struct {
	Port    int
	Timeout time.Duration
}

func Load(ctx context.Context, configFilePath string) (*Config, error) {
	yamlFile, err := cfg.YAMLFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var (
		serverPort = cfg.Setting[int]{
			Default:   func() int { return 8080 },
			Validator: cfg.Between(5000, 9000),
			Providers: []cfg.Provider[int]{
				cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
				yamlFile.Int("server", "port"),
			},
		}
		serverTimeout = cfg.Setting[time.Duration]{
			Validator: cfg.Between(0, time.Minute),
			Providers: []cfg.Provider[time.Duration]{
				cfg.EnvVar("APP_SERVER_TIMEOUT", time.ParseDuration),
				yamlFile.Duration("server", "timeout"),
			},
		}
		dbHost = cfg.Setting[string]{
			Default: func() string { return "localhost" },
			Providers: []cfg.Provider[string]{
				cfg.EnvVarStr("APP_DB_HOST"),
				yamlFile.String("database", "host"),
			},
		}
		dbPort = cfg.Setting[int]{
			Default: func() int { return 3306 },
			Providers: []cfg.Provider[int]{
				cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
				yamlFile.Int("database", "port"),
			},
		}
		dbUsername = cfg.Setting[string]{
			Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
			Providers: []cfg.Provider[string]{
				cfg.EnvVarStr("APP_DB_USERNAME"),
				yamlFile.String("database", "username"),
			},
		}
		dbPassword = cfg.Setting[string]{
			Providers: []cfg.Provider[string]{
				cfg.EnvVarStr("APP_DB_PASSWORD"),
				yamlFile.String("database", "password"),
			},
		}
	)

	// Call ValidateAll to ensure the subsequent MustGet calls won't panic.
	if err := cfg.ValidateAll(
		ctx,
		serverPort,
		serverTimeout,
		dbHost,
		dbPort,
		dbUsername,
		dbPassword); err != nil {
		return nil, err
	}

	c := &Config{
		Server: ServerConfig{
			Port:    serverPort.MustGet(ctx),
			Timeout: serverTimeout.MustGet(ctx),
		},
		DB: DBConfig{
			Host:     dbHost.MustGet(ctx),
			Port:     dbPort.MustGet(ctx),
			Username: dbUsername.MustGet(ctx),
			Password: dbPassword.MustGet(ctx),
		},
	}

	return c, nil
}
