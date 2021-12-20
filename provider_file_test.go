package cfg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileProviderJSON(t *testing.T) {
	p, err := File(FormatJSON, "testdata/config.json")
	assert.NoError(t, err)

	port := p.Provide("server", "port")
	out, err := port.Provide(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, string(out), "8000")
}

func TestFileProviderYAML(t *testing.T) {
	p, err := File(FormatYAML, "testdata/config.yaml")
	assert.NoError(t, err)

	port := p.Provide("server", "port")
	out, err := port.Provide(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, string(out), "8000")
}
