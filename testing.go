package cfg

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/zpatrick/testx/assert"
)

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

func AssertProvides[T comparable](t testing.TB, p Provider[T], expected T) {
	out, err := p.Provide(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, out, expected)
}
