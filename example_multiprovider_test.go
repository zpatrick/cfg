package cfg_test

import (
	"context"
	"fmt"
	"os"

	"github.com/zpatrick/cfg"
)

type MyConfig struct {
	EnableTLS bool

	DBHost string
	DBPort int
}

func ExampleMultiProvider() {
	defaultYaml, err := cfg.YAMLFile("config.default.yaml")
	if err != nil {
		panic(err)
	}

	// config.local.env, config.test.env, config.production.env, ...
	environmentYaml, err := cfg.YAMLFile(fmt.Sprintf("config.%s.yaml", os.Getenv("APP_ENV")))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	errs := cfg.NewErrorAggregator()

	c := MyConfig{
		EnableTLS: cfg.Setting[bool]{
			Name: "EnableTLS",
			Provider: cfg.MultiProvider[bool]{
				defaultYaml.Bool("enable_tls"),
				environmentYaml.Bool("enable_tls"),
			},
		}.MustGet(ctx, errs),
		DBHost: cfg.Setting[string]{
			Name: "DBHost",
			Provider: cfg.MultiProvider[string]{
				defaultYaml.String("database", "host"),
				environmentYaml.String("database", "host"),
			},
		}.MustGet(ctx, errs),
		DBPort: cfg.Setting[int]{
			Name:    "DBPort",
			Default: cfg.Pointer(3306),
			Provider: cfg.MultiProvider[int]{
				defaultYaml.Int("database", "port"),
				environmentYaml.Int("database", "port"),
			},
		}.MustGet(ctx, errs),
	}

	if err := errs.Err(); err != nil {
		panic(err)
	}

	fmt.Println(c)
}
