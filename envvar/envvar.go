package envvar

import (
	"context"
	"os"

	"github.com/zpatrick/cfg"
)

// New returns a string provider for the given environment variable.
func New(key string) cfg.Provider[string] {
	return &provider[string]{
		key:    key,
		decode: func(s string) (string, error) { return s, nil },
	}
}

// Newf returns a provider for the given environment variable.
// The decode function is used to convert the raw environment variable into type T.
func Newf[T any](key string, decode func(string) (T, error)) cfg.Provider[T] {
	return &provider[T]{key: key, decode: decode}
}

type provider[T any] struct {
	key    string
	decode func(string) (T, error)
}

func (p *provider[T]) Provide(ctx context.Context) (out T, err error) {
	val := os.Getenv(p.key)
	if val == "" {
		return out, cfg.NoValueProvidedError
	}

	return p.decode(val)
}
