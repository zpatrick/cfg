package cfg_test

import (
	"context"
	"strings"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func TestLoader_includesKeyInErrorMessage(t *testing.T) {
	err := cfg.Load(context.Background(), cfg.Schemas{
		"foo": cfg.Schema[int]{},
	})

	assert.Error(t, err)
	assert.Equal(t, strings.Contains(err.Error(), "foo"), true)
}
