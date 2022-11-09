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

func ExampleEnvVar() {
	const key = "APP_PORT"
	os.Setenv(key, "9090")

	appPortProvider := cfg.EnvVar(key, strconv.Atoi)
	appPort, err := appPortProvider.Provide(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("port is:", appPort)
	// Output: port is: 9090
}

func TestEnvVarNoValue(t *testing.T) {
	const key = "APP_PORT"
	os.Unsetenv(key)

	provider := cfg.EnvVar(key, strconv.Atoi)
	_, err := provider.Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
