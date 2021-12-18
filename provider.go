package cfg

import "context"

// A Provider is used to load a "raw" configuration value from some predetermined source.
// A Decoder will be used to convert the value from byte slice into the appropriate type.
type Provider interface {
	Provide(ctx context.Context) ([]byte, error)
}

// The ProviderFunc is an adapter type which allows ordinary functions to be used as Providers.
type ProviderFunc func(ctx context.Context) ([]byte, error)

// Provide calls f(ctx).
func (f ProviderFunc) Provide(ctx context.Context) ([]byte, error) {
	return f(ctx)
}
