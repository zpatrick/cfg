package cfg

import (
	"context"

	"github.com/pkg/errors"
)

// Pointer returns a pointer to t. This can be used to assign a Setting.Default in one line.
func Pointer[T any](t T) *T {
	return &t
}

// Get calls s.Get, returning the provided value and adding the error (if present) to errs.
func Get[T any](ctx context.Context, errs *ErrorAggregator, s Setting[T]) T {
	out, err := s.Get(ctx)
	if err != nil {
		errs.Add(errors.Wrapf(err, "failed to get %s", s.Name))
	}

	return out
}

type Setting[T any] struct {
	// Name is an optional field which provides a human-friendly identifier for a setting.
	// Currently, this is only used for adding details to error messages.
	Name string
	// Provider is a required field which provides a value for the setting.
	Provider Provider[T]
	// Default is an optional field which specifies the default value to use if
	// no value is provided by Provider.
	Default *T
	// Validator is an optional field which ensures the value provided is valid.
	Validator Validator[T]
}

// Get returns the value provided by s.Provider.
// If s.Validator is set, it will be used to validate the provided value.
// If no value is provided and s.Default is set, *s.Default will be returned.
// If no value is provided and s.Default is not set, a NoValueProvidedError will be returned.
func (s Setting[T]) Get(ctx context.Context) (T, error) {
	val, err := s.Provider.Provide(ctx)
	if err != nil {
		if errors.Is(err, NoValueProvidedError) {
			if s.Default != nil {
				return *s.Default, nil
			}

			return val, NoValueProvidedError
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
