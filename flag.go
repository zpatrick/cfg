package cfg

import (
	"context"
	"flag"
)

// Flag returns a provider from the given flag's pointer.
// The set and name values are used to check if the flag was explicitly set or not.
// If the flag is not explicitly set, a NoValueProvidedError will be returned.
func Flag[T any](ptr *T, set *flag.FlagSet, name string) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		var flagIsSet bool
		set.Visit(func(f *flag.Flag) {
			if f.Name == name {
				flagIsSet = true
			}
		})

		if !flagIsSet {
			return out, NoValueProvidedError
		}

		return *ptr, nil
	})
}

// FlagWithDefault returns a provider from the given flag's pointer.
// The flag's default value will be returned if the flag is not explicitly set.
func FlagWithDefault[T any](ptr *T) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (T, error) {
		return *ptr, nil
	})
}
