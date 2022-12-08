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

// Val returns the value provided by s.Load.
// This method will panic if s.Load has not been called or if the the method returned an error.
func (s Setting[T]) Val() T {
	return *s.val
}

// ValOK peforms the same logic as Val, but returns a boolean instead of panicing
// if there is not a valid value to return.
func (s Setting[T]) ValOK() (T, bool) {
	if s.val == nil {
		var t T
		return t, false
	}

	return *s.val, true
}
