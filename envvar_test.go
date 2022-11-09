package cfg_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

type Config struct {
	Host string
	Port int
}

func ExampleEnvVar() {
	ctx := context.Background()
	errs := cfg.NewErrorAggregator()

	c := Config{
		Host: cfg.Setting[string]{
			Name:     "host",
			Provider: cfg.EnvVarStr("APP_HOST"),
		}.MustGet(ctx, errs),
		Port: cfg.Setting[int]{
			Name:     "port",
			Provider: cfg.EnvVar("APP_PORT", strconv.Atoi),
		}.MustGet(ctx, errs),
	}

	if err := errs.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("host: %s, port: %d\n", c.Host, c.Port)
}

func TestEnvVarNoValue(t *testing.T) {
	const key = "APP_PORT"
	os.Unsetenv(key)

	provider := cfg.EnvVar(key, strconv.Atoi)
	_, err := provider.Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
