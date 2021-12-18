package cfg

import (
	"context"
	"os"
)

// EnvVar returns a provider for the given environment variable.
func EnvVar(key string) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		val := os.Getenv(key)
		if val == "" {
			return nil, NoValueProvidedError
		}

		return []byte(val), nil
	})
}
