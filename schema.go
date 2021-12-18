
package cfg

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"context"
)

var NoValueError = errors.New("no value set")

type Provider interface {
	Provide(ctx context.Context) ([]byte, error)
}

type ProviderFunc func(ctx context.Context) ([]byte, error)

func (f ProviderFunc) Provide(ctx context.Context) ([]byte, error) {
	return f(ctx)
}

type Schema[T any] struct {
	Name      string
	Decode    func(b []byte) (T, error)
	Default   func() T
	Validate  func(T) error
	Providers []Provider
}

func (s Schema[T]) Load(ctx context.Context) (out T, err error) {
	for _, p := range s.Providers {
		b, err := p.Provide(ctx)
		if err != nil {
			if errors.Is(err, NoValueError) {
				continue
			}

			return out, fmt.Errorf("failed to load %s: %w", s.Name, err)
		}

		out, err = s.Decode(b)
		if err != nil {
			return out, fmt.Errorf("failed to decode %s: %s", s.Name, err)
		}

		return out, nil
	}

	if s.Default != nil {
		return s.Default(), nil
	}

	return out, NoValueError
}

func DecodeInt(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}



func EnvVar(key string) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		if val := os.Getenv(key); val != "" {
			return []byte(val), nil
		}

		return nil, NoValueError
	})
}

func OneOf[T comparable](allowed ...T) func(in T) error {
	s := make(map[T]struct{}, len(allowed))
	for _, a := range allowed {
		s[a] = struct{}{}
	}

	return func(in T) error {
		if _, ok := s[in]; ok {
			return nil
		}

		return fmt.Errorf("value %v is not in %v", in, allowed)
	}
}

// TODO: use interface, anything with 'Load' method
// func Validate(ss ...Schema) error {
// 	for _, s := range ss {
// 		if err := s.Load(); err != nil {
// 			return fmt.Errorf("failed to load %s: %v")
// 		}
// 	}
// }

// func main() {
// 	os.Setenv("APP_PORT", "9096")
//  file := cfg.NewFileProvider(cfg.FormatYAML, "config.yaml")

// 	port := Schema[int]{
// 		Name:     "port",
// 		Decode:   DecodeInt,
// 		Default:  func() int { return 9090 },
// 		Validate: Between(1, 5),
// 		Providers: []Provider{
// 			EnvVar("APP_PORT"),
// 			file.Provide("server", "port"),
// 		},
// 	}

// if err := Validate(port); err != nil { panic(err) }

// 	fmt.Println(port.Load())
// }
