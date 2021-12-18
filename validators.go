package cfg

import (
	"constraints"
	"fmt"
)

// Between returns a validator which ensures the input is > min and < max.
func Between[T constraints.Ordered](min, max T) func(input T) error {
	return func(input T) error {
		switch {
		case input <= min:
			return fmt.Errorf("input %v is <= than the allowed minimum %v", input, min)
		case input >= max:
			return fmt.Errorf("input %v is >= than the allowed maximum %v", input, max)
		default:
			return nil
		}
	}
}

// Contains returns a validator which ensures the input is equal to one of the given vals.
func Contains[T comparable](vals ...T) func(input T) error {
	set := make(map[T]struct{}, len(vals))
	for _, v := range vals {
		set[v] = struct{}{}
	}

	return func(input T) error {
		if _, ok := set[input]; ok {
			return nil
		}

		return fmt.Errorf("input %v not contained in %v", input, vals)
	}
}