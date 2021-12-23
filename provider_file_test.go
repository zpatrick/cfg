package cfg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileProviderINI(t *testing.T) {
	f, err := File(FormatINI, "testdata/config.ini")
	if err != nil {
		t.Fatal(err)
	}

	testFileProviderHelper(t, f)
}

func TestFileProviderJSON(t *testing.T) {
	f, err := File(FormatJSON, "testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	testFileProviderHelper(t, f)
}

func TestFileProviderYAML(t *testing.T) {
	f, err := File(FormatYAML, "testdata/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	testFileProviderHelper(t, f)
}

func testFileProviderHelper(t *testing.T, f *fileProvider) {
	assert.Equal(t, "8000", mustProvide(t, f, "server", "port"))
	assert.Equal(t, "30s", mustProvide(t, f, "server", "request_timeout"))
	assert.Equal(t, "false", mustProvide(t, f, "server", "enable_ssl"))

	assert.Equal(t, "localhost", mustProvide(t, f, "database", "host"))
	assert.Equal(t, "3306", mustProvide(t, f, "database", "port"))
	assert.Equal(t, "root", mustProvide(t, f, "database", "username"))
	assert.Equal(t, "secret", mustProvide(t, f, "database", "password"))
}

func mustProvide(t *testing.T, f *fileProvider, section string, keys ...string) string {
	out, err := f.Provide(section, keys...).Provide(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return DecodeString(out)
}
