package cfg_test

import (
	"context"
	"fmt"
	"os"

	"github.com/zpatrick/cfg"
)

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

func ExampleProvider_custom() {
	var env Environment
	schema := cfg.Schema[Environment]{
		Dest:      &env,
		Default:   cfg.Addr(Development),
		Validator: cfg.OneOf(Development, Staging, Production),
		Provider: cfg.ProviderFunc[Environment](func(context.Context) (Environment, error) {
			appEnv := os.Getenv("APP_ENV")
			if appEnv == "" {
				return "", cfg.NoValueProvidedError
			}

			return Environment(appEnv), nil
		}),
	}

	os.Setenv("APP_ENV", "staging")
	if err := schema.Load(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println(env)
	// Output: staging
}
