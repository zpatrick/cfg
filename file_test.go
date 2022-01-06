package cfg_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func TestFileProviderINI(t *testing.T) {
	f, err := cfg.File(cfg.ParseINI(), "testdata/config.ini")
	assert.NilError(t, err)

	mustProvide(t, 8000, cfg.Convert(strconv.Atoi, f.String("server", "port")))
	mustProvide(t, time.Second*30, cfg.Convert(time.ParseDuration, f.String("server", "request_timeout")))
	mustProvide(t, true, cfg.Convert(strconv.ParseBool, f.String("server", "enable_ssl")))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, cfg.Convert(strconv.Atoi, f.String("database", "port")))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func TestFileProviderJSON(t *testing.T) {
	f, err := cfg.File(cfg.ParseJSON(), "testdata/config.json")
	assert.NilError(t, err)

	mustProvide(t, 8000, cfg.Convert(cfg.Float64ToInt, f.Float64("server", "port")))
	mustProvide(t, time.Second*30, cfg.Convert(cfg.StringToDuration, f.String("server", "request_timeout")))
	mustProvide(t, true, f.Bool("server", "enable_ssl"))

	mustProvide(t, "localhost", f.String("database", "host"))
	mustProvide(t, 3306, cfg.Convert(cfg.Float64ToInt, f.Float64("database", "port")))
	mustProvide(t, "root", f.String("database", "username"))
	mustProvide(t, "secret", f.String("database", "password"))
}

func TestFileProviderYAML(t *testing.T) {
	f, err := cfg.File(cfg.ParseYAML(), "testdata/config.yaml")
	assert.NilError(t, err)

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
