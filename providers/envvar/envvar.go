package envvar

import (
	"context"
	"os"

	"github.com/zpatrick/cfg"
)

func New(key string) cfg.Provider[string] {
	return &envVarProvider[string]{
		key:    key,
		decode: func(s string) (string, error) { return s, nil },
	}
}

// EnvVar returns a provider for the given environment variable.
// The decode function is used to convert the raw environment variable into type T.
func Newf[T any](key string, decode func(string) (T, error)) cfg.Provider[T] {
	return &envVarProvider[T]{key: key, decode: decode}
}

type envVarProvider[T any] struct {
	key    string
	decode func(string) (T, error)
}

func (e *envVarProvider[T]) Provide(ctx context.Context) (out T, err error) {
	val := os.Getenv(e.key)
	if val == "" {
		return out, cfg.NoValueProvidedError
	}

	return e.decode(val)
}
