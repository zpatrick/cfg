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

const (
	Development = "development"
	Staging     = "staging"
	Production  = "production"
)

func Load(ctx context.Context, configFilePath string) (*Config, error) {
	env, err := cfg.Setting[string]{
		Name:      "environment",
		Default:   cfg.Pointer(Development),
		Validator: cfg.OneOf(Development, Staging, Production),
		Provider:  cfg.EnvVarStr("APP_ENV"),
	}.Get(ctx)
	if err != nil {
		return nil, err
	}

	yamlFile, err := cfg.YAMLFile(configFilePath)
	if err != nil {
		return nil, err
	}

	errs := cfg.NewErrorAggregator()
	c := Config{
		Server: server.Config{
			EnableTLS: env == Production,
			Port: cfg.Setting[int]{
				Name:      "Server Port",
				Default:   cfg.Pointer(8080),
				Validator: cfg.Between(5000, 9000),
				Provider: cfg.MultiProvider[int]{
					cfg.EnvVar("APP_SERVER_PORT", strconv.Atoi),
					yamlFile.Int("server", "port"),
				},
			}.MustGet(ctx, errs),
			Timeout: cfg.Setting[time.Duration]{
				Name:      "Server Timeout",
				Validator: cfg.Between(0, time.Minute),
				Provider: cfg.MultiProvider[time.Duration]{
					cfg.EnvVar("APP_SERVER_TIMEOUT", time.ParseDuration),
					yamlFile.Duration("server", "timeout"),
				},
			}.MustGet(ctx, errs),
		},
		DB: database.Config{
			Host: cfg.Setting[string]{
				Name:    "DB Host",
				Default: cfg.Pointer("localhost"),
				Provider: cfg.MultiProvider[string]{
					cfg.EnvVarStr("APP_DB_HOST"),
					yamlFile.String("database", "host"),
				},
			}.MustGet(ctx, errs),
			Port: cfg.Setting[int]{
				Name:    "DB Port",
				Default: cfg.Pointer(3306),
				Provider: cfg.MultiProvider[int]{
					cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
					yamlFile.Int("database", "port"),
				},
			}.MustGet(ctx, errs),
			Username: cfg.Setting[string]{
				Name:      "DB Username",
				Validator: cfg.OneOf("admin", "app_rw", "app_ro"),
				Provider: cfg.MultiProvider[string]{
					cfg.EnvVarStr("APP_DB_USERNAME"),
					yamlFile.String("database", "username"),
				},
			}.MustGet(ctx, errs),
			Password: cfg.Setting[string]{
				Name: "DB Password",
				Provider: cfg.MultiProvider[string]{
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
