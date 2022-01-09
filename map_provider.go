package cfg

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type mapProvider map[string]any

// Traverse ...
func (m mapProvider) Traverse(key string, keys ...string) (any, error) {
	val, ok := m[key]
	if !ok {
		return nil, NoValueProvidedError
	}

	if len(keys) == 0 {
		return val, nil
	}

	switch val := val.(type) {
	case map[string]any:
		return mapProvider(val).Traverse(keys[0], keys[1:]...)
	case map[any]any:
		return asMapProvider(val).Traverse(keys[0], keys[1:]...)
	default:
		uerr := NewUnexpectedTypeError(map[string]any{}, val)
		return nil, errors.Wrapf(uerr, "unable to traverse past key %s", key)
	}
}

// asMapProvider extracts all key/val pairs from m where key is of type string.
func asMapProvider(m map[any]any) mapProvider {
	out := make(map[string]any, len(m))
	for key, val := range m {
		if key, ok := key.(string); ok {
			out[key] = val
		}
	}

	return out
}

func (m mapProvider) String(section string, keys ...string) Provider[string] {
	return mapProvide[string](m, section, keys...)
}

func (m mapProvider) Int(section string, keys ...string) Provider[int] {
	return mapProvide[int](m, section, keys...)
}

func (m mapProvider) Bool(section string, keys ...string) Provider[bool] {
	return mapProvide[bool](m, section, keys...)
}

func (m mapProvider) Duration(section string, keys ...string) Provider[time.Duration] {
	return ProviderFunc[time.Duration](func(ctx context.Context) (out time.Duration, err error) {
		val, err := mapProvide[string](m, section, keys...).Provide(ctx)
		if err != nil {
			return out, err
		}

		return time.ParseDuration(val)
	})
}

func mapProvide[T any](m mapProvider, section string, keys ...string) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		val, err := m.Traverse(section, keys...)
		if err != nil {
			return out, NoValueProvidedError
		}

		v, ok := val.(T)
		if !ok {
			return out, NewUnexpectedTypeError(out, val)
		}

		return v, nil
	})
}
