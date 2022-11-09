package cfg

import (
	"context"

	"github.com/pkg/errors"
)

type Setting[T any] struct {
	Name string
	// Default specifies the default value to use if no value is provided by Providers.
	// If this field is nil, a NoValueProvidedError will be returned instead.
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

			return val, errors.Wrapf(err, "%s: failed to load", s.Name)
		}

		if s.Validator != nil {
			if err := s.Validator.Validate(val); err != nil {
				return val, errors.Wrapf(err, "%s: failed to validate", s.Name)
			}
		}

		return val, nil
	}

	if s.Default != nil {
		return s.Default(), nil
	}

	return val, errors.Wrapf(NoValueProvidedError, s.Name)
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
