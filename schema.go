package cfg

import (
	"context"
	"errors"
	"fmt"
)

// A Schema ...
type Schema[T any] struct {
	Name      string
	Default   func() T
	Validator Validator[T]
	Providers []Provider[T]
}

func (s Schema[T]) Load(ctx context.Context) (val T, err error) {
	for _, p := range s.Providers {
		val, err := p.Provide(ctx)
		if err != nil {
			if errors.Is(err, NoValueProvidedError) {
				continue
			}

			return val, fmt.Errorf("failed to load %s: %w", s.Name, err)
		}

		if s.Validator != nil {
			if err := s.Validator.Validate(val); err != nil {
				return val, err
			}
		}

		return val, nil
	}

	if s.Default != nil {
		return s.Default(), nil
	}

	return val, NoValueProvidedError
}

func (s Schema[T]) MustLoad(ctx context.Context) T {
	out, err := s.Load(ctx)
	if err != nil {
		panic(err)
	}

	return out
}

func (s Schema[T]) Validate(ctx context.Context) error {
	_, err := s.Load(ctx)
	return err
}
