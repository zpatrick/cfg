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

	addr := &cfg.Setting[string]{
		Provider: envvar.New("APP_ADDR"),
	}

	addr.Load(context.Background())
	fmt.Println(addr.Val())
	// Output: localhost
}

func ExampleNewf() {
	os.Setenv("APP_PORT", "8080")

	port := &cfg.Setting[int]{
		Provider: envvar.Newf("APP_PORT", strconv.Atoi),
	}

	port.Load(context.Background())
	fmt.Println(port.Val())
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
