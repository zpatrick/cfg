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

// TODO:
// func eample() {
// c, err := cfg.Load(ctx)
// if err != nil {
//   return nil, err
// }

// svr := http.Server{
//   Port: c.Port.Value,
//   Host: c.Host.Value,
// }

// // config.go
// type Config struct {
//   Port cfg.Setting[int]
//   Host cfg.Setting[string]
// }

// func Load(ctx context.Context) (*Config, error) {
//   yamlFile, err := yaml.NewFile(ctx, FilePath)
//   if err != nil {
//     return nil, err
//   }

//   c := &Config{
//     Port: cfg.Setting[int]{
//       Name: "port",
//       Default: cfg.Pointer(9090),
//       Validate: cfg.Between(0, 9999),
//       Provider: cfg.EnvVar("APP_PORT"),
//     },
//     Host: cfg.Setting[string]{
//       Name: "host",
//       Default: cfg.Pointer("localhost"),
//       Provider: cfg.MultiProvider{
//         cfg.EnvVar("APP_HOST"),
//         yamlFile.String("server", "host"),
//       },
//     },
//   }

//   if err := cfg.Load(ctx, cfg); err != nil {
//     return nil, err
//   }

//   return c, nil
// }

// }

func Load(ctx context.Context, c any) error {
	// ptr := reflect.ValueOf(c)
	// if ptr.Type().Kind() != reflect.Ptr {
	// 	return fmt.Errorf("parameter is not a pointer to a struct")
	// }

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

		setting, ok := field.Interface().(interface{ Load(context.Context) error })
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

	//  if valueOf.Type().Kind() == reflect.Ptr {
	//  if valueOf.Type().Kind() == reflect.Struct {
	//val.Type().Field(0).Name

	return nil
}
