// Package cfg allows developers to define complex configuration for their applications with minimal code.
// This package has been designed to help satisfy the needs of teams who are building microservice in go.
// The goals of this package include:
//   - Allow teams to use consistent patterns for configuration across different applications.
//   - Coalesce multiple sources of configuration.
//   - Custom validation of configuration values.
//   - House a variety of tools to work with common configuration sources/formats.
package cfg

import (
	"context"

	"github.com/pkg/errors"
)

// Loader is the interface implemented by types that can load values into themselves.
type Loader interface {
	Load(context.Context) error
}

func Load(ctx context.Context, schemas map[string]Loader) error {
	for name, schema := range schemas {
		if err := schema.Load(ctx); err != nil {
			return errors.Wrapf(err, "failed to load %s", name)
		}
	}

	return nil
}

// Pointer returns a pointer to t. This can be used to assign a Schema.Default in one line.
func Pointer[T any](t T) *T {
	return &t
}
