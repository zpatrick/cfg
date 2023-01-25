package cfg

import (
	"context"

	"github.com/pkg/errors"
)

// Schemas is a named type for a map of Schemas.
type Schemas map[string]interface {
	Load(context.Context) error
	schema()
}

// A Schema models a configuration setting for an application.
type Schema[T any] struct {
	// Dest is a required field which points where to store the configuration value when Load is called.
	Dest *T
	// Provider is a required field which provides the configuration value.
	Provider Provider[T]
	// Default is an optional field which specifies the fallback value to use if Provider returns a NoValueProvidedError.
	Default *T
	// Validator is an optional field which ensures the value provided by Provider is valid.
	Validator Validator[T]
}

func (Schema[T]) schema() {}

// Load calls s.Provider.Provide to get the configuration value and store it into s.Dest.
// If s.Provider.Provide returns a NoValueProvidedError and s.Default is not nil, s.Default will be used instead.
// If s.Validator is set, it will be used to validate the provided value.
func (s Schema[T]) Load(ctx context.Context) error {
	if s.Dest == nil {
		return errors.New("dest is nil")
	}

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

	*s.Dest = val
	return nil
}
