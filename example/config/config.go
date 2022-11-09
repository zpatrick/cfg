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

func Load(ctx context.Context, configFilePath string) (*Config, error) {
	yamlFile, err := cfg.YAMLFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var errs cfg.ErrorAggregator
	c := Config{
		Server: server.Config{
			Port: cfg.Setting[int]{
				Name:      "Server Port",
				Default:   func() int { return 8080 },
				Validator: cfg.Between(5000, 9000),
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
					yamlFile.Int("server", "port"),
				},
			}.MustGet(ctx, errs),
			Timeout: cfg.Setting[time.Duration]{
				Name:      "Server Timeout",
				Validator: cfg.Between(0, time.Minute),
				Providers: []cfg.Provider[time.Duration]{
					cfg.EnvVar("APP_SERVER_TIMEOUT", time.ParseDuration),
					yamlFile.Duration("server", "timeout"),
				},
			}.MustGet(ctx, errs),
		},
		DB: database.Config{
			Host: cfg.Setting[string]{
				Name:    "DB Host",
				Default: func() string { return "localhost" },
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_HOST"),
					yamlFile.String("database", "host"),
				},
			}.MustGet(ctx, errs),
			Port: cfg.Setting[int]{
				Name:    "DB Port",
				Default: func() int { return 3306 },
				Providers: []cfg.Provider[int]{
					cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
					yamlFile.Int("database", "port"),
				},
			}.MustGet(ctx, errs),
			Username: cfg.Setting[string]{
				Name:      "DB Username",
				Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_USERNAME"),
					yamlFile.String("database", "username"),
				},
			}.MustGet(ctx, errs),
			Password: cfg.Setting[string]{
				Name: "DB Password",
				Providers: []cfg.Provider[string]{
					cfg.EnvVarStr("APP_DB_PASSWORD"),
					yamlFile.String("database", "password"),
				},
			}.MustGet(ctx, errs),
		},
	}

	if err := errs.Err(); err != nil {
		return nil, err
	}

	return &c, nil
}
