package cfg_test

import (
	"context"
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

type TestConfig struct {
	private int
	A       *cfg.Setting[int]
	B       *cfg.Setting[int]
	Nested  NestedTestConfig
}

type NestedTestConfig struct {
	C *cfg.Setting[int]
	D *cfg.Setting[int]
}

func TestLoad(t *testing.T) {
	c := &TestConfig{
		private: 1,
		A: &cfg.Setting[int]{
			Provider: cfg.StaticProvider(1),
		},
		B: &cfg.Setting[int]{
			Provider: cfg.StaticProvider(2),
		},
		Nested: NestedTestConfig{
			C: &cfg.Setting[int]{
				Provider: cfg.StaticProvider(3),
			},
			D: &cfg.Setting[int]{
				Provider: cfg.StaticProvider(4),
			},
		},
	}

	if err := cfg.Load(context.Background(), c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, c.A.Val(), 1)
	assert.Equal(t, c.B.Val(), 2)
	assert.Equal(t, c.Nested.C.Val(), 3)
	assert.Equal(t, c.Nested.D.Val(), 4)
}
