package cfg_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

type embedded struct {
	E1 cfg.Setting[int]
	E2 *cfg.Setting[int]
}

type nested struct {
	N1 cfg.Setting[int]
	N2 *cfg.Setting[int]
}

func TestLoad(t *testing.T) {
	c := struct {
		embedded
		// TODO: star embedded
		Nested    nested
		NestedPtr *nested
		P1        cfg.Setting[int]
		P2        *cfg.Setting[int]
		priv      cfg.Setting[int]
		Interface io.Writer
		Literal   *http.Request
	}{
		embedded: embedded{
			E1: cfg.Setting[int]{Provider: cfg.StaticProvider(1)},
			E2: &cfg.Setting[int]{Provider: cfg.StaticProvider(2)},
		},
		Nested: nested{
			N1: cfg.Setting[int]{Provider: cfg.StaticProvider(3)},
			N2: &cfg.Setting[int]{Provider: cfg.StaticProvider(4)},
		},
		NestedPtr: &nested{
			N1: cfg.Setting[int]{Provider: cfg.StaticProvider(5)},
			N2: &cfg.Setting[int]{Provider: cfg.StaticProvider(6)},
		},
		P1:        cfg.Setting[int]{Provider: cfg.StaticProvider(7)},
		P2:        &cfg.Setting[int]{Provider: cfg.StaticProvider(8)},
		priv:      cfg.Setting[int]{},
		Interface: io.Discard,
		Literal:   &http.Request{},
	}

	if err := cfg.Load(context.Background(), &c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, c.embedded.E1.Val(), 1)
	assert.Equal(t, c.embedded.E2.Val(), 2)
	assert.Equal(t, c.Nested.N1.Val(), 3)
	assert.Equal(t, c.Nested.N2.Val(), 4)
	assert.Equal(t, c.NestedPtr.N1.Val(), 5)
	assert.Equal(t, c.NestedPtr.N2.Val(), 6)
	assert.Equal(t, c.P1.Val(), 7)
	assert.Equal(t, c.P2.Val(), 8)
}
