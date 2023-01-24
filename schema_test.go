package cfg_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func ExampleSchema_validation() {
	var userName string
	schema := cfg.Schema[string]{
		Dest:      &userName,
		Validator: cfg.OneOf("admin", "guest"),
		Provider:  cfg.StaticProvider("other"),
	}

	if err := schema.Load(context.Background()); err != nil {
		fmt.Println(err)
	}

	// Output: validation failed: input other not contained in [admin guest]
}

func TestSchemaLoad_populatesDest(t *testing.T) {
	var out int
	port := cfg.Schema[int]{
		Dest:     &out,
		Provider: cfg.StaticProvider(8080),
	}

	assert.NilError(t, port.Load(context.Background()))
	assert.Equal(t, out, 8080)
}

func TestSchemaLoad_returnsUnhandledProviderError(t *testing.T) {
	var out int
	port := cfg.Schema[int]{
		Dest:    &out,
		Default: cfg.Addr(8080),
		Provider: cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, io.EOF
		}),
	}

	assert.ErrorIs(t, port.Load(context.Background()), io.EOF)
}

func TestSchemaLoad_usesDefaultWhenHandlingNoValueProvidedError(t *testing.T) {
	var out int
	port := cfg.Schema[int]{
		Dest:    &out,
		Default: cfg.Addr(8080),
		Provider: cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, cfg.NoValueProvidedError
		}),
	}

	assert.NilError(t, port.Load(context.Background()))
	assert.Equal(t, out, 8080)
}

func TestSchemaLoad_returnsValidationError(t *testing.T) {
	var (
		out    int
		called bool
	)

	port := cfg.Schema[int]{
		Dest:     &out,
		Provider: cfg.StaticProvider(8080),
		Validator: cfg.ValidatorFunc[int](func(i int) error {
			called = true
			return io.EOF
		}),
	}

	assert.ErrorIs(t, port.Load(context.Background()), io.EOF)
	assert.Equal(t, called, true)
}

func TestSchemaLoad_validationSuccess(t *testing.T) {
	var (
		out    int
		called bool
	)

	port := cfg.Schema[int]{
		Dest:     &out,
		Provider: cfg.StaticProvider(8080),
		Validator: cfg.ValidatorFunc[int](func(i int) error {
			called = true
			return nil
		}),
	}

	assert.NilError(t, port.Load(context.Background()))
	assert.Equal(t, called, true)
}

func TestSchemaSatisfiesLoaderInterface(t *testing.T) {
	var _ cfg.Loader = cfg.Schema[int]{}
}
