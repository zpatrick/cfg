package ini_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal"
	"github.com/zpatrick/cfg/providers/ini"
	"github.com/zpatrick/testx/assert"
)

func TestNew(t *testing.T) {
	const data = `
root = "hello"

[server]
enabled = true
port = 8080
percent = 8.5
timeout = "5s"
`

	path, err := internal.WriteTempFile(t.TempDir(), data)
	assert.NilError(t, err)
	t.Cleanup(func() { os.Remove(path) })

	f, err := ini.New(path)
	assert.NilError(t, err)

	internal.AssertProvides(t, f.String("", "root"), "hello")
	internal.AssertProvides(t, f.Bool("server", "enabled"), true)
	internal.AssertProvides(t, f.Int64("server", "port"), 8080)
	internal.AssertProvides(t, f.Float64("server", "percent"), 8.5)
	internal.AssertProvides(t, f.Duration("server", "timeout"), time.Second*5)

	_, err = f.Int("invalid", "port").Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)

	_, err = f.Int("server", "invalid").Provide(context.Background())
	assert.ErrorIs(t, err, cfg.NoValueProvidedError)
}
