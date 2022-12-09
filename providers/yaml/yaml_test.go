package yaml_test

import (
	"os"
	"testing"

	"github.com/zpatrick/cfg"
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

	path, err := cfg.WriteTempFile(t.TempDir(), data)
	assert.NilError(t, err)
	t.Cleanup(func() { os.Remove(path) })

	f, err := yaml.New(path)
	assert.NilError(t, err)

	//cfg.AssertProvides(t, f.String("root"), "hello")
	//cfg.AssertProvides(t, f.Int("server", "port"), 8080)
	cfg.AssertProvides(t, f.Bool("server", "enabled"), true)
	//cfg.AssertProvides(t, f.Duration("server", "timeout", "read"), time.Second*5)
	//cfg.AssertProvides(t, f.Float64("server", "percent"), 8.5)

	// _, err = f.Int("invalid").Provide(context.Background())
	// assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
