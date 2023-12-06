package cfg

import (
	"context"
	"reflect"

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
// If allowZero is false, the Provider will return a NoValueProvidedError if v is the zero value.
//
//	var p Provider[int] = StaticProvider(5)
func StaticProvider[T any](v T, allowZero bool) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (T, error) {
		if !allowZero && reflect.ValueOf(v).IsZero() {
			return v, NoValueProvidedError
		}

		return v, nil
	})
}

// StaticProviderAddr adapts pv into a Provider[T].
// If pv is nil, the Provider will return a NoValueProvidedError.
// If allowZero is false, the Provider will return a NoValueProvidedError if v is the zero value.
//
//	var p Provider[int] = StaticProviderAddr(&addr)
func StaticProviderAddr[T any](pv *T, allowZero bool) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (T, error) {
		var zero T
		if pv == nil {
			return zero, NoValueProvidedError
		}

		if !allowZero && reflect.ValueOf(*pv).IsZero() {
			return zero, NoValueProvidedError
		}

		return *pv, nil
	})
}
