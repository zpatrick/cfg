package cfg

import (
	"fmt"

	"golang.org/x/exp/constraints"

	"go.uber.org/multierr"
)

// A Validator checks whether or not a given value T is considered valid.
type Validator[T any] interface {
	// Validate returns an error if input is considered invalid.
	Validate(input T) error
}

// A ValidatorFunc is an adapter type which allows functions to be used as Validators.
type ValidatorFunc[T any] func(T) error

// Validate calls f(in).
func (f ValidatorFunc[T]) Validate(in T) error {
	return f(in)
}

// Between returns a validator which ensures the input is >= min and <= max.
func Between[T constraints.Ordered](min, max T) Validator[T] {
	return ValidatorFunc[T](func(in T) error {
		switch {
		case in < min:
			return fmt.Errorf("input %v is < than the allowed minimum %v", in, min)
		case in > max:
			return fmt.Errorf("input %v is > than the allowed maximum %v", in, max)
		default:
			return nil
		}
	})
}

// OneOf returns a validator which ensures the input is equal to one of the given vals.
func OneOf[T comparable](vals ...T) Validator[T] {
	set := make(map[T]struct{}, len(vals))
	for _, v := range vals {
		set[v] = struct{}{}
	}

	return ValidatorFunc[T](func(in T) error {
		if _, ok := set[in]; ok {
			return nil
		}

		return fmt.Errorf("input %v not contained in %v", in, vals)
	})
}

// And combines the given validators into a single validator,
// requiring each validator check to succeed.
func And[T any](a, b Validator[T], vs ...Validator[T]) Validator[T] {
	all := append([]Validator[T]{a, b}, vs...)
	return ValidatorFunc[T](func(in T) error {
		for _, v := range all {
			if err := v.Validate(in); err != nil {
				return err
			}
		}

		return nil
	})
}

// Or combines the given validators into a single validator,
// requiring only one validator check to succeed.
func Or[T any](a, b Validator[T], vs ...Validator[T]) Validator[T] {
	all := append([]Validator[T]{a, b}, vs...)
	return ValidatorFunc[T](func(in T) error {
		var errs []error
		for _, v := range all {
			if err := v.Validate(in); err != nil {
				errs = append(errs, err)
				continue
			}

			return nil
		}

		return multierr.Combine(errs...)
	})
}
