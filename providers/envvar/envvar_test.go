package envvar_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/envvar"
	"github.com/zpatrick/testx/assert"
)

func ExampleNew() {
	var addr string
	schema := &cfg.Schema[string]{
		Dest:     &addr,
		Provider: envvar.New("APP_ADDR"),
	}

	os.Setenv("APP_ADDR", "localhost")
	if err := schema.Load(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println(addr)
	// Output: localhost
}

func ExampleNewf() {
	var port int
	schema := &cfg.Schema[int]{
		Dest:     &port,
		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
	}

	os.Setenv("APP_PORT", "9090")
	if err := schema.Load(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println(port)
	// Output: 9090
}

func TestNew_returnsProperValue(t *testing.T) {
	const key = "TEST_KEY"
	os.Setenv(key, "value")

	val, err := envvar.New(key).Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, val, "value")
}

func TestNew_returnsNoValueProvidedError(t *testing.T) {
	os.Clearenv()

	_, err := envvar.New("TEST_KEY").Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}

func TestNewf_returnsDecodedValue(t *testing.T) {
	const key = "TEST_KEY"
	os.Setenv(key, "value")

	decoder := func(s string) (int, error) {
		assert.Equal(t, s, "value")
		return 1, nil
	}

	val, err := envvar.Newf("TEST_KEY", decoder).Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, val, 1)
}

func TestNewf_returnsNoValueProvidedError(t *testing.T) {
	os.Clearenv()

	_, err := envvar.Newf("TEST_KEY", strconv.Atoi).Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
