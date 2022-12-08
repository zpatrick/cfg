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
	os.Setenv("APP_ADDR", "localhost")

	addr := cfg.Setting[string]{
		Name:     "address",
		Provider: envvar.New("APP_ADDR"),
	}

	val, _ := addr.Get(context.Background())
	fmt.Println(val)
	// Output: localhost
}

func ExampleNewf() {
	os.Setenv("APP_PORT", "8080")

	port := cfg.Setting[int]{
		Name:     "port",
		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
	}

	val, _ := port.Get(context.Background())
	fmt.Println(val)
	// Output: 8080
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
