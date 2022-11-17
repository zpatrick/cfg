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
	env := cfg.Setting[Environment]{
		Name:      "environment",
		Default:   cfg.Pointer(Development),
		Validator: cfg.OneOf(Development, Staging, Production),
		Provider: cfg.ProviderFunc[Environment](func(context.Context) (Environment, error) {
			appEnv := os.Getenv("APP_ENV")
			if appEnv == "" {
				return "", cfg.NoValueProvidedError
			}

			return Environment(appEnv), nil
		}),
	}

	val, _ := env.Get(context.Background())
	fmt.Println(val)
	// Output: development
}
