package cfg

import (
	"context"

	"github.com/pkg/errors"
)

type Setting[T any] struct {
	Name string
	// Default specifies the default value to use if no value is provided by the Provider.
	// If this field is nil, a NoValueProvidedError will be returned instead.
	Default   *T
	Validator Validator[T]
	Provider  Provider[T]
}

func (s Setting[T]) Get(ctx context.Context) (T, error) {
	val, err := s.Provider.Provide(ctx)
	if err != nil {
		if errors.Is(err, NoValueProvidedError) {
			if s.Default != nil {
				return *s.Default, nil
			}

			return val, err
		}

		return val, err
	}

	if s.Validator != nil {
		if err := s.Validator.Validate(val); err != nil {
			return val, errors.Wrapf(err, "validation failed")
		}
	}

	return val, nil
}

func (s Setting[T]) MustGet(ctx context.Context, errs *ErrorAggregator) T {
	out, err := s.Get(ctx)
	if err != nil {
		errs.Add(errors.Wrapf(err, "failed to get %s", s.Name))
	}

	return out
}
