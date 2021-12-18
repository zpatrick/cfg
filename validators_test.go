package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetween(t *testing.T) {
	assert.NoError(t, Between(0, 2)(1))
	assert.NoError(t, Between(-1, 1)(0))
	assert.NoError(t, Between(8.5, 8.6)(8.55))
}

func TestBetweenError(t *testing.T) {
	assert.Error(t, Between(0, 2)(0))
	assert.Error(t, Between(0, 2)(2))
	assert.Error(t, Between(0, 2)(-1))
	assert.Error(t, Between(0, 2)(3))
	assert.Error(t, Between(8.5, 8.6)(8.4))
}

func TestContains(t *testing.T) {
	assert.NoError(t, Contains(0)(0))
	assert.NoError(t, Contains(0, 1)(1))
	assert.NoError(t, Contains(0, 1, 2)(1))
	assert.NoError(t, Contains(5, 5, 5)(5))
}

func TestContainsError(t *testing.T) {
	assert.Error(t, Contains[int]()(0))
	assert.Error(t, Contains(0)(1))
	assert.Error(t, Contains(0, 1, 2)(3))
}
