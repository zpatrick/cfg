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

// Loader is the interface implemented by types that can load values into themselves.
type Loader interface {
	Load(context.Context) error
}

// Load takes a struct, c, and traverses the fields recursively.
// If any field is of type Setting, the Load method will be called on that field.
// Additionally, the Load method will be called on any field which implements the Loader interface.
func Load(ctx context.Context, c any) error {
	val := reflect.Indirect(reflect.ValueOf(c))
	if val.Type().Kind() != reflect.Struct {
		return fmt.Errorf("parameter is not of type struct")
	}

	for _, e := range reflect.VisibleFields(val.Type()) {
		if !e.IsExported() {
			continue
		}

		field := val.FieldByIndex(e.Index)
		loader, ok := checkIfLoader(field)
		if ok {
			if err := loader.Load(ctx); err != nil {
				return errors.Wrapf(err, "failed to load field %s", e.Name)
			}

			continue
		}

		if reflect.Indirect(field).Kind() == reflect.Struct {
			if field.Kind() != reflect.Ptr && field.CanAddr() {
				field = field.Addr()
			}

			if err := Load(ctx, field.Interface()); err != nil {
				return errors.Wrapf(err, "failed to load field %s", e.Name)
			}
		}
	}

	return nil
}

func checkIfLoader(v reflect.Value) (Loader, bool) {
	if v.CanInterface() {
		loader, ok := v.Interface().(Loader)
		if ok {
			return loader, true
		}
	}

	if v.CanAddr() {
		return checkIfLoader(v.Addr())
	}

	return nil, false
}

// Are you a struct?
// For each of your fields

// func Load(ctx context.Context, c any) error {
// 	val := reflect.Indirect(reflect.ValueOf(c))
// 	if val.Type().Kind() != reflect.Struct {
// 		return fmt.Errorf("parameter is not a struct")
// 	}

// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		fieldName := val.Type().Field(i).Name
// 		print(fieldName)

// 		if !field.CanInterface() {
// 			continue
// 		}

// 		//field.Addr().MethodByName("Load").Call([]reflect.Value{reflect.ValueOf(ctx)})
// 		//setting, ok := field.Interface().(Loader)
// 		// TODO: check if can interface above
// 		// if field.Kind() == reflect.Ptr {
// 		// 	field = field.Elem()
// 		// }
// 		field = reflect.Indirect(field)

// 		//v := field.Interface()

// 		// TODO: check canAddr
// 		// OOOOH SHIT OK THIS IS IT!
// 		vptr := field.Addr()

// 		//_, ok := field.Elem().Interface().(Loader)
// 		_, ok := vptr.Interface().(Loader)
// 		if ok {
// 			// TODO: check output err
// 			field.Addr().MethodByName("Load").Call([]reflect.Value{reflect.ValueOf(ctx)})

// 			// if err := setting.Load(ctx); err != nil {
// 			// 	return errors.Wrapf(err, "failed to load field %s", fieldName)
// 			// }

// 			continue
// 		}

// 		// if field.Type().Kind() == reflect.Struct {
// 		// 	if err := Load(ctx, field.Interface()); err != nil {
// 		// 		return errors.Wrapf(err, "failed to load field %s", fieldName)
// 		// 	}
// 		// }
// 	}

// 	return nil
// }

// func Load2(ctx context.Context, c any) error {
// 	ptr := reflect.ValueOf(c)
// 	if ptr.Kind() != reflect.Ptr {
// 		return errors.New("c must be a pointer to a struct")
// 	}

// 	// reflect.Indirect(ptr)
// 	val := ptr.Elem()
// 	if val.Kind() != reflect.Struct {
// 		return errors.New("c must be a pointer to a struct")
// 	}

// 	for _, field := range reflect.VisibleFields(val.Type()) {
// 		if !field.IsExported() {
// 			continue
// 		}

// 		var loader Loader
// 		t := reflect.TypeOf(&loader).Elem()
// 		if field.Type.Implements(reflect.TypeOf(t)) {
// 			field := val.FieldByName(field.Name)
// 			field.MethodByName("Loader").Call([]reflect.Value{reflect.ValueOf(ctx)})
// 		}

// 	}

// 	// for i := 0; i < val.NumField(); i++ {
// 	// 	field := val.Field(i)
// 	// 	if !field.CanInterface() {
// 	// 		continue
// 	// 	}

// 	// 	field.Kind()
// 	// 	loader, ok := field.Interface().(Loader)
// 	// 	if ok {
// 	// 		if err := loader.Load(ctx); err != nil {
// 	// 			return err
// 	// 		}

// 	// 		continue
// 	// 	}
// 	//}

// 	return nil
// }
