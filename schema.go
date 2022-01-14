package cfg

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type Setting[T any] struct {
	Default   func() T
	Validator Validator[T]
	Providers []Provider[T]
}

func (s *Setting[T]) Get(ctx context.Context) (val T, err error) {
	for _, p := range s.Providers {
		val, err := p.Provide(ctx)
		if err != nil {
			if errors.Is(err, NoValueProvidedError) {
				continue
			}

			return val, errors.Wrapf(err, "failed to load setting")
		}

		if s.Validator != nil {
			if err := s.Validator.Validate(val); err != nil {
				return val, errors.Wrapf(err, "failed to validate setting")
			}
		}

		return val, nil
	}

	if s.Default != nil {
		return s.Default(), nil
	}

	return val, NoValueProvidedError
}

func (s *Setting[T]) MustGet(ctx context.Context) T {
	out, err := s.Get(ctx)
	if err != nil {
		panic(err)
	}

	return out
}

func (s Setting[T]) Validate(ctx context.Context) error {
	_, err := s.Get(ctx)
	return err
}

type Settings map[string]Validateable

func (s Settings) Validate(ctx context.Context) error {
	errs := []error{}
	for key, v := range s {
		if err := v.Validate(ctx); err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to validate %s", key))
		}
	}

	return multierr.Combine(errs...)
}

func Get[T any](ctx context.Context, setting any) (out T, err error) {
	s, ok := setting.(Setting[T])
	if !ok {
		return out, fmt.Errorf("setting is not of type Setting[%s]", reflect.TypeOf(out).String())
	}

	return s.Get(ctx)
}

func MustGet[T any](ctx context.Context, setting any) (out T) {
	out, err := Get[T](ctx, setting)
	if err != nil {
		panic(err)
	}

	return out
}
