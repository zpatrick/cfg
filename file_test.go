package cfg_test

import (
	"context"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func TestFileProviderINI(t *testing.T) {
	t.Skip()
	f, err := cfg.File(cfg.ParseINI(), "testdata/config.ini")
	assert.NilError(t, err)

	testFileProviderHelper(t, f)
}

func TestFileProviderJSON(t *testing.T) {
	f, err := cfg.File(cfg.ParseJSON(), "testdata/config.json")
	assert.NilError(t, err)

	testFileProviderHelper(t, f)
}

func TestFileProviderYAML(t *testing.T) {
	f, err := cfg.File(cfg.ParseYAML(), "testdata/config.yaml")
	assert.NilError(t, err)

	testFileProviderHelper(t, f)
}

func testFileProviderHelper(t *testing.T, f *cfg.FileProvider) {
	mustProvide(t, 8000, cfg.Convert(cfg.Float64ToInt, f.Float64("server", "port")))
	mustProvide(t, time.Second*30, cfg.Convert(cfg.StringToDuration, f.String("server", "request_timeout")))
	mustProvide(t, true, f.Bool("server", "enable_ssl"))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, cfg.Convert(cfg.Float64ToInt, f.Float64("database", "port")))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func mustProvide[T comparable](t testing.TB, expected T, p cfg.Provider[T]) {
	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, expected)
}
