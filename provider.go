package cfg

import (
	"context"

	"github.com/pkg/errors"
)

// A Provider loads a configuration value from a predetermined source.
// If no value is provided by the underlying source, the Provider must return a NoValueProvidedError.
type Provider[T any] interface {
	Provide(ctx context.Context) (T, error)
}

// The ProviderFunc is an adapter type which allows ordinary functions to be used as Providers.
type ProviderFunc[T any] func(ctx context.Context) (T, error)

// Provide calls f(ctx).
func (f ProviderFunc[T]) Provide(ctx context.Context) (T, error) {
	return f(ctx)
}

// MultiProvider allows a slice of Provider[T] to be used as a Provider[T].
type MultiProvider[T any] []Provider[T]

// Provide iterates through each provider in m and returns the first value given by a provider.
// If a provider returns a NoValueProvidedError, the iteration continues.
// If no providers return a value, a NoValueProvidedError is returned.
func (m MultiProvider[T]) Provide(ctx context.Context) (T, error) {
	var t T
	for _, p := range m {
		val, err := p.Provide(ctx)
		if err != nil {
			if errors.Is(err, NoValueProvidedError) {
				continue
			}

			return t, err
		}

		return val, nil
	}

	return t, NoValueProvidedError
}

// StaticProvider adapts v into a Provider[T].
//
//	var p Provider[int] = StaticProvider(5)
func StaticProvider[T any](v T) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (T, error) {
		return v, nil
	})
}
