package cfg_test

import (
	"context"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func TestFileProviderINI(t *testing.T) {
	f, err := cfg.INIFile("testdata/config.ini")
	assert.NilError(t, err)

	mustProvide(t, 8000, f.Int("server", "port"))
	mustProvide(t, time.Second*30, f.Duration("server", "request_timeout"))
	mustProvide(t, true, f.Bool("server", "enable_ssl"))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, f.Int("database", "port"))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func TestFileProviderJSON(t *testing.T) {
	f, err := cfg.JSONFile("testdata/config.json")
	assert.NilError(t, err)

	mustProvide(t, 8000, cfg.Convert(cfg.Float64ToInt, f.Float64("server", "port")))
	mustProvide(t, time.Second*30, f.Duration("server", "request_timeout"))
	mustProvide(t, true, f.Bool("server", "enable_ssl"))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, cfg.Convert(cfg.Float64ToInt, f.Float64("database", "port")))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func TestFileProviderTOML(t *testing.T) {
	f, err := cfg.TOMLFile("testdata/config.toml")
	assert.NilError(t, err)

	mustProvide(t, 8000, f.Int64("server", "port"))
	mustProvide(t, time.Second*30, f.Duration("server", "request_timeout"))
	mustProvide(t, true, f.Bool("server", "enable_ssl"))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, f.Int64("database", "port"))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func TestFileProviderYAML(t *testing.T) {
	f, err := cfg.YAMLFile("testdata/config.yaml")
	assert.NilError(t, err)

	mustProvide(t, 8000, f.Int("server", "port"))
	mustProvide(t, time.Second*30, f.Duration("server", "request_timeout"))
	mustProvide(t, true, f.Bool("server", "enable_ssl"))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, f.Int("database", "port"))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func mustProvide[T comparable](t testing.TB, expected T, p cfg.Provider[T]) {
	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, expected)
}
