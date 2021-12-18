package cfg

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvVar(t *testing.T) {
	const key = "CFG_TEST_KEY"
	os.Setenv(key, "value")

	out, err := EnvVar(key).Provide(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, string(out), "value")
}

func TestEnvVarNoValue(t *testing.T) {
	const key = "CFG_TEST_KEY"
	os.Unsetenv(key)

	_, err := EnvVar(key).Provide(context.Background())
	assert.ErrorIs(t, err, NoValueProvidedError)
}
