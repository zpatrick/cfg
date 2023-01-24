package toml_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal/cfgtest"
	"github.com/zpatrick/cfg/providers/toml"
	"github.com/zpatrick/testx/assert"
)

func TestNew(t *testing.T) {
	const data = `
root = "hello"

[servers]
enabled = true

  [servers.alpha]
  port = 8080
  percent = 8.5
  timeout = "5s"
`

	path, err := cfgtest.WriteTempFile(t.TempDir(), data)
	assert.NilError(t, err)
	t.Cleanup(func() { os.Remove(path) })

	f, err := toml.New(path)
	assert.NilError(t, err)

	cfgtest.AssertProvides(t, f.String("root"), "hello")
	cfgtest.AssertProvides(t, f.Bool("servers", "enabled"), true)
	cfgtest.AssertProvides(t, f.Int64("servers", "alpha", "port"), 8080)
	cfgtest.AssertProvides(t, f.Float64("servers", "alpha", "percent"), 8.5)
	cfgtest.AssertProvides(t, f.Duration("servers", "alpha", "timeout"), time.Second*5)

	_, err = f.Int("invalid").Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
