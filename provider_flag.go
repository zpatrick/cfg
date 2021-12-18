package cfg

import (
	"context"
	"flag"
	"strconv"
	"time"
)

type flagDefaultChecker func() error

func IgnoreFlagDefault(name string, visitFunc func(fn func(*flag.Flag))) flagDefaultChecker {
	return func() error {
		var flagIsSet bool
		visitFunc(func(f *flag.Flag) {
			if f.Name == name {
				flagIsSet = true
			}
		})

		if !flagIsSet {
			return NoValueProvidedError
		}

		return nil
	}
}

// IntFlag returns a provider from the given flag's pointer.
func IntFlag(ptr *int, checker flagDefaultChecker) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		if checker != nil {
			if err := checker(); err != nil {
				return nil, err
			}
		}

		s := strconv.Itoa(*ptr)
		return []byte(s), nil
	})
}

// StringFlag returns a provider from the given flag's pointer.
func StringFlag(ptr *string, checker flagDefaultChecker) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		if checker != nil {
			if err := checker(); err != nil {
				return nil, err
			}
		}

		return []byte(*ptr), nil
	})
}

// BoolFlag returns a provider from the given flag's pointer.
func BoolFlag(ptr *bool, checker flagDefaultChecker) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		if checker != nil {
			if err := checker(); err != nil {
				return nil, err
			}
		}

		s := strconv.FormatBool(*ptr)
		return []byte(s), nil
	})
}

func DurationFlag(ptr *time.Duration, checker flagDefaultChecker) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		if checker != nil {
			if err := checker(); err != nil {
				return nil, err
			}
		}

		return []byte(ptr.String()), nil
	})
}

// TODO: 64s, uints
