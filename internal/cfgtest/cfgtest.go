package cfgtest

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

// WriteTempFile creates a new temporary file under dir with the contents of data.
// The file name is returned.
func WriteTempFile(dir, data string) (string, error) {
	file, err := ioutil.TempFile(dir, "config.yaml")
	if err != nil {
		return "", nil
	}

	if _, err := bytes.NewBufferString(data).WriteTo(file); err != nil {
		return "", err
	}

	return file.Name(), nil
}

// AssertProvides calls t.Fatal if p.Provide doesn't return expected.
func AssertProvides[T comparable](t testing.TB, p cfg.Provider[T], expected T) {
	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, expected)
}
