package cfg

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvVar(t *testing.T) {
	const key = "CFG_TEST_KEY"
	os.Setenv(key, "5")

	provider := EnvVar(key, strconv.Atoi)
	out, err := provider.Provide(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, out, 5)
}

func TestEnvVarNoValue(t *testing.T) {
	const key = "CFG_TEST_KEY"
	os.Unsetenv(key)

	provider := EnvVar(key, strconv.Atoi)
	_, err := provider.Provide(context.Background())
	assert.ErrorIs(t, err, NoValueProvidedError)
}
