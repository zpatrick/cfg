package cfg_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/envvar"
	"github.com/zpatrick/testx/assert"
)

func ExampleSchema_multiProvider() {
	var userName string
	schema := cfg.Schema[string]{
		Dest: &userName,
		// Note that order matters when using MultiProvider:
		// We first will use USERNAME_ALPHA if that envvar is set,
		// falling back to using USERNAME_BRAVO if not.
		Provider: cfg.MultiProvider[string]{
			envvar.New("USERNAME_ALPHA"),
			envvar.New("USERNAME_BRAVO"),
		},
	}

	os.Setenv("USERNAME_ALPHA", "foo")
	os.Setenv("USERNAME_BRAVO", "bar")

	if err := schema.Load(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println(userName)

	// Output: foo
}

func TestMultiProvider_returnsFirstError(t *testing.T) {
	p := cfg.MultiProvider[int]{
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, io.EOF
		}),
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, nil
		}),
	}

	_, err := p.Provide(context.Background())
	assert.ErrorIs(t, err, io.EOF)
}

func TestMultiProvider_returnsFirstValue(t *testing.T) {
	p := cfg.MultiProvider[int]{
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 1, nil
		}),
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, errors.New("we shouldn't have gotten this far")
		}),
	}

	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, 1)
}

func TestMultiProvider_iteratesThroughNoValueProvidedError(t *testing.T) {
	p := cfg.MultiProvider[int]{
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 1, cfg.NoValueProvidedError
		}),
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 2, nil
		}),
	}

	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, 2)
}

func TestMultiProvider_returnsNoValueProvidedErrorWhenDoneIterating(t *testing.T) {
	p := cfg.MultiProvider[int]{
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, cfg.NoValueProvidedError
		}),
		cfg.ProviderFunc[int](func(context.Context) (int, error) {
			return 0, cfg.NoValueProvidedError
		}),
	}

	_, err := p.Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}

func TestStaticProvider(t *testing.T) {
	p := cfg.StaticProvider(5, false)
	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, 5)
}
