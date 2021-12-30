package cfg

import (
	"context"
	"errors"
	"fmt"
)

// A Schema ...
type Schema[T any] struct {
	Name      string
	Decode    func(b []byte) (T, error)
	Default   func() T
	Validator Validator[T]
	Providers []Provider
}

func (s Schema[T]) Load(ctx context.Context) (out T, err error) {
	for _, p := range s.Providers {
		b, err := p.Provide(ctx)
		if err != nil {
			if errors.Is(err, NoValueProvidedError) {
				continue
			}

			return out, fmt.Errorf("failed to load %s: %w", s.Name, err)
		}

		out, err = s.Decode(b)
		if err != nil {
			return out, fmt.Errorf("failed to decode %s: %s", s.Name, err)
		}

		if s.Validator != nil {
			if err := s.Validator.Validate(out); err != nil {
				return out, err
			}
		}

		return out, nil
	}

	if s.Default != nil {
		return s.Default(), nil
	}

	return out, NoValueProvidedError
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
