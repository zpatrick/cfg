package envvar_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/envvar"
	"github.com/zpatrick/testx/assert"
)

func TestEnvVarNoValue(t *testing.T) {
	const key = "APP_PORT"
	os.Unsetenv(key)

	provider := envvar.Newf(key, strconv.Atoi)
	_, err := provider.Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
