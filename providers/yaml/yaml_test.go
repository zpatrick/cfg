package yaml_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/yaml"
	"github.com/zpatrick/testx/assert"
)

func writeTempFile(dir, data string) (string, error) {
	file, err := ioutil.TempFile(dir, "config.yaml")
	if err != nil {
		return "", nil
	}

	if _, err := bytes.NewBufferString(data).WriteTo(file); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func assertProvides[T comparable](t testing.TB, p cfg.Provider[T], expected T) {
	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, expected)
}

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

	path, err := writeTempFile(t.TempDir(), data)
	assert.NilError(t, err)
	t.Cleanup(func() { os.Remove(path) })

	f, err := yaml.New(path)
	assert.NilError(t, err)

	assertProvides(t, f.String("root"), "hello")
	assertProvides(t, f.Int("server", "port"), 8080)
	assertProvides(t, f.Bool("server", "enabled"), true)
	assertProvides(t, f.Duration("server", "timeout", "read"), time.Second*5)
	assertProvides(t, f.Float64("server", "percent"), 8.5)
}
