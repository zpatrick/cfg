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
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// Loader is the interface implemented by types that can load setting values into themselves.
type Loader interface {
	Load(context.Context) error
}

// Load takes a struct, c, and traverses the fields recursively.
// If any field is of type Setting, the Load method will be called on that field.
// Additionally, the Load method will be called on any field which implements the Loader interface.
func Load(ctx context.Context, c any) error {
	val := reflect.Indirect(reflect.ValueOf(c))
	if val.Type().Kind() != reflect.Struct {
		return fmt.Errorf("parameter is not a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		if !field.CanInterface() {
			continue
		}

		setting, ok := field.Interface().(Loader)
		if ok {
			if err := setting.Load(ctx); err != nil {
				return errors.Wrapf(err, "failed to load field %s", fieldName)
			}

			continue
		}

		if field.Type().Kind() == reflect.Struct {
			if err := Load(ctx, field.Interface()); err != nil {
				return errors.Wrapf(err, "failed to load field %s", fieldName)
			}
		}
	}

	return nil
}
