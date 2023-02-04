package yaml_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal"
	"github.com/zpatrick/cfg/providers/yaml"
	"github.com/zpatrick/testx/assert"
)

func TestNew(t *testing.T) {
	const data = `
root: hello
server:
  port: 8080
  percent: 8.5
  enabled: true
  timeout:
    read: 5s
`

	path, err := internal.WriteTempFile(t.TempDir(), data)
	assert.NilError(t, err)
	t.Cleanup(func() { os.Remove(path) })

	f, err := yaml.New(path)
	assert.NilError(t, err)

	internal.AssertProvides(t, f.String("root"), "hello")
	internal.AssertProvides(t, f.Int("server", "port"), 8080)
	internal.AssertProvides(t, f.Bool("server", "enabled"), true)
	internal.AssertProvides(t, f.Duration("server", "timeout", "read"), time.Second*5)
	internal.AssertProvides(t, f.Float64("server", "percent"), 8.5)

	_, err = f.Int("invalid").Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)

	_, err = f.Int("server", "invalid").Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
