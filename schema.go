package cfg

import (
	"context"

	"github.com/pkg/errors"
)

/*
Let's take some hard stances here:

1. Passing around a config using builtin types is superior to cfg-defined types.
	a) Consumer's code is decoupled from the package; can swap out as necessary.
	b) The cfg-defined type adds no additional functionality (just calls .Val)
	c) The config package is a small part of an application, we shouldn't optimize for that package's types.


2. Clear is better than clever
	a) Using reflect and ctx jamming is a nice api, but a bit unclear what is happening.
	b) I


3. A clean api sells
	a) I don't want to require any more than 1 struct for configuration declaration.


2. API Ideas:
	a) Take the flag package idea of returning pointers to values, then let Load deal with the errors later.
	b) cfg.Load top-level function.
	c) setting.Load bottom-level function.







type MyConfig struct {
	Host string
	Port int
	Username string
}



func Load(ctx context.Context) {
	var c MyConfig

	if err := cfg.Load(ctx, map[string]cfg.Loader{
		"port": cfg.Schema[int]{
			Dest: &c.Port,
			Default: cfg.Pointer(8080),
			Validate: cfg.Between(5000, 9000),
			Provider: envvar.Newf("APP_PORT", strconv.Itoa),
		},
		"host": cfg.Schema[string]{
			Dest: &c.Host,
			Default: cfg.Pointer("localhost"),
			Provider: envvar.Newf("APP_HOST"),
		},
		"db.username": cfg.Schema[string]{
			Dest: &c.DBUsername,
			Provider: cfg.MultiProvider{
				envvar.Newf("APP_DB_USERNAME"),
				yamlFile.String("db", "username"),
			},
		},
	})


	return c, nil
}
*/

// Pointer returns a pointer to t. This can be used to assign a Schema.Default in one line.
func Pointer[T any](t T) *T {
	return &t
}

type Schema[T any] struct {
	// Dest is a required field that points to a T in which to store the configuration value.
	Dest *T
	// Provider is a required field which provides a value for the setting.
	Provider Provider[T]
	// Default is an optional field which specifies the default value to use if
	// no value is provided by Provider.
	Default *T
	// Validator is an optional field which ensures the value provided is valid.
	Validator Validator[T]

	val *T
}

// Get loads the value provided by s.Provider.
// If s.Validator is set, it will be used to validate the provided value.
// If no value is provided and s.Default is set, s.Default will be used.
// If no value is provided and s.Default is not set, a NoValueProvidedError will be returned.
func (s *Schema[T]) Load(ctx context.Context) error {
	if s.Dest == nil {
		return errors.New("dest is nil")
	}

	if s.Provider == nil {
		return errors.New("provider is nil")
	}

	val, err := s.Provider.Provide(ctx)
	if err != nil {
		if !errors.Is(err, NoValueProvidedError) {
			return err
		}

		if s.Default == nil {
			return err
		}

		val = *s.Default
	}

	if s.Validator != nil {
		if err := s.Validator.Validate(val); err != nil {
			return errors.Wrapf(err, "validation failed")
		}
	}

	s.val = &val
	return nil
}
