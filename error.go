package cfg

import (
	"go.uber.org/multierr"
)

type sentinelError string

func (s sentinelError) Error() string {
	return string(s)
}

const NoValueProvidedError sentinelError = "no value provided"

type ErrorAggregator struct {
	errors []error
}

func NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{
		errors: []error{},
	}
}

func (e *ErrorAggregator) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *ErrorAggregator) Err() error {
	return multierr.Combine(e.errors...)
}
