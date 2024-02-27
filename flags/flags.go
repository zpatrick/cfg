package flags

import (
	"context"
	"flag"

	"github.com/zpatrick/cfg"
)

// New returns a provider from the given flag's pointer.
// The set and name values are used to check if the flag was explicitly set or not.
// If the flag is not explicitly set, a NoValueProvidedError will be returned.
func New[T any](set *flag.FlagSet, ptr *T, name string) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		var flagIsSet bool
		set.Visit(func(f *flag.Flag) {
			if f.Name == name {
				flagIsSet = true
			}
		})

		if !flagIsSet {
			return out, cfg.NoValueProvidedError
		}

		return *ptr, nil
	})
}

// NewWithDefault returns a provider from the given flag's pointer.
// The flag's default value will be returned if the flag is not explicitly set.
func NewWithDefault[T any](ptr *T) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (T, error) {
		return *ptr, nil
	})
}
