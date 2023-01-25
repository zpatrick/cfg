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

// Load calls Load(ctx) for each schema in schemas.
func Load(ctx context.Context, schemas Schemas) error {
	for name, schema := range schemas {
		if err := schema.Load(ctx); err != nil {
			return errors.Wrapf(err, "failed to load %s", name)
		}
	}

	return nil
}

// Addr returns the address of t.
func Addr[T any](t T) *T {
	return &t
}
