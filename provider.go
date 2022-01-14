package cfg

import (
	"context"
)

// A Provider loads a configuration value from some predetermined source.
// If no value is provided by the underlying source, the Provider must return
// a NoValueProvidedError.
type Provider[T any] interface {
	Provide(ctx context.Context) (T, error)
}

// The ProviderFunc is an adapter type which allows ordinary functions to be used as Providers.
type ProviderFunc[T any] func(ctx context.Context) (T, error)

// Provide calls f(ctx).
func (f ProviderFunc[T]) Provide(ctx context.Context) (T, error) {
	return f(ctx)
}

// This adapts an instance of a type T turn itself into a Provider[T].
//   var p Provider[int] = Providable[int]{5}
type Providable[T any] struct{ V T }

// Provide returns p.V.
func (p Providable[T]) Provide(context.Context) (T, error) { return p.V, nil }
