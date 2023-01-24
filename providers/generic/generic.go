package generic

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/zpatrick/cfg"
)

type UnexpectedTypeError struct {
	Expected, Provided reflect.Type
}

func NewUnexpectedTypeError(expectedVal, providedVal interface{}) *UnexpectedTypeError {
	return &UnexpectedTypeError{
		Expected: reflect.TypeOf(expectedVal),
		Provided: reflect.TypeOf(providedVal),
	}
}

func (e *UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type error: provided type was %s (expected %s)", e.Provided, e.Expected)
}

type Provider map[string]any

func (p Provider) Get(key string, keys ...string) (any, error) {
	val, ok := p[key]
	if !ok {
		return nil, cfg.NoValueProvidedError
	}

	if len(keys) == 0 {
		return val, nil
	}

	switch val := val.(type) {
	case map[string]any:
		return Provider(val).Get(keys[0], keys[1:]...)
	case map[any]any:
		return asProvider(val).Get(keys[0], keys[1:]...)
	default:
		uerr := NewUnexpectedTypeError(map[string]any{}, val)
		return nil, errors.Wrapf(uerr, "unable to traverse past key %s", key)
	}
}

// asProvider extracts all key/val pairs from m where key is of type string.
func asProvider(m map[any]any) Provider {
	out := make(map[string]any, len(m))
	for key, val := range m {
		if key, ok := key.(string); ok {
			out[key] = val
		}
	}

	return out
}

func (p Provider) String(section string, keys ...string) cfg.Provider[string] {
	return provide[string](p, section, keys...)
}

func (p Provider) Float64(section string, keys ...string) cfg.Provider[float64] {
	return provide[float64](p, section, keys...)
}

func (p Provider) Int(section string, keys ...string) cfg.Provider[int] {
	return provide[int](p, section, keys...)
}

func (p Provider) Int64(section string, keys ...string) cfg.Provider[int64] {
	return provide[int64](p, section, keys...)
}

func (p Provider) Uint64(section string, keys ...string) cfg.Provider[uint64] {
	return provide[uint64](p, section, keys...)
}

func (p Provider) Bool(section string, keys ...string) cfg.Provider[bool] {
	return provide[bool](p, section, keys...)
}

func (p Provider) Duration(section string, keys ...string) cfg.Provider[time.Duration] {
	return cfg.ProviderFunc[time.Duration](func(ctx context.Context) (out time.Duration, err error) {
		val, err := provide[string](p, section, keys...).Provide(ctx)
		if err != nil {
			return out, err
		}

		return time.ParseDuration(val)
	})
}

func provide[T any](p Provider, section string, keys ...string) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		val, err := p.Get(section, keys...)
		if err != nil {
			return out, cfg.NoValueProvidedError
		}

		v, ok := val.(T)
		if !ok {
			return out, NewUnexpectedTypeError(out, val)
		}

		return v, nil
	})
}
