package cfg

import (
	"context"

	"github.com/pkg/errors"
)

// Pointer returns a pointer to t. This can be used to assign a Setting.Default in one line.
func Pointer[T any](t T) *T {
	return &t
}

type Setting[T any] struct {
	// Provider is a required field which provides a value for the setting.
	Provider Provider[T]
	// Default is an optional field which specifies the default value to use if
	// no value is provided by Provider.
	Default *T
	// Validator is an optional field which ensures the value provided is valid.
	Validator Validator[T]

	val *T
}

// Get loads the value provided by s.Provider.
// If s.Validator is set, it will be used to validate the provided value.
// If no value is provided and s.Default is set, s.Default will be used.
// If no value is provided and s.Default is not set, a NoValueProvidedError will be returned.
func (s *Setting[T]) Load(ctx context.Context) error {
	if s.Provider == nil {
		return errors.New("provider is nil")
	}

	val, err := s.Provider.Provide(ctx)
	if err != nil {
		if !errors.Is(err, NoValueProvidedError) {
			return err
		}

		if s.Default == nil {
			return err
		}

		val = *s.Default
	}

	if s.Validator != nil {
		if err := s.Validator.Validate(val); err != nil {
			return errors.Wrapf(err, "validation failed")
		}
	}

	s.val = &val
	return nil
}

// Val returns the value loaded into s using s.Load.
// If s.Load has not been called, and s.Default is set, the default will be returned.
// Otherwise, the zero value of T will be returned.
func (s Setting[T]) Val() T {
	if s.val != nil {
		return *s.val
	}

	if s.Default != nil {
		return *s.Default
	}

	var zero T
	return zero
}
