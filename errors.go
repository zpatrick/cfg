package cfg

import (
	"fmt"
	"reflect"
)

type sentinelError string

func (s sentinelError) Error() string {
	return string(s)
}

const NoValueProvidedError sentinelError = "no value provided"

type UnexpectedTypeError struct {
	Expected, Provided reflect.Type
}

func NewUnexpectedTypeError(expectedVal, providedVal interface{}) *UnexpectedTypeError {
	return &UnexpectedTypeError{
		Expected: reflect.TypeOf(expectedVal),
		Provided: reflect.TypeOf(providedVal),
	}
}

func (e *UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type error: provided type was %s (expected %s)", e.Provided, e.Expected)
}
