package cfg

import (
	"testing"

	"github.com/zpatrick/testx/assert"
)

func TestBetween(t *testing.T) {
	assert.NilError(t, Between(0, 2).Validate(0))
	assert.NilError(t, Between(0, 2).Validate(1))
	assert.NilError(t, Between(0, 2).Validate(2))
	assert.NilError(t, Between(8.5, 8.6).Validate(8.55))
}

func TestBetweenError(t *testing.T) {
	assert.Error(t, Between(0, 2).Validate(-1))
	assert.Error(t, Between(0, 2).Validate(3))
	assert.Error(t, Between(8.5, 8.6).Validate(8.4))
	assert.Error(t, Between(8.5, 8.6).Validate(8.7))
}

func TestOneOf(t *testing.T) {
	assert.NilError(t, OneOf(0).Validate(0))
	assert.NilError(t, OneOf(0, 1).Validate(1))
	assert.NilError(t, OneOf(0, 1, 2).Validate(1))
	assert.NilError(t, OneOf(5, 5, 5).Validate(5))
}

func TestOneOfError(t *testing.T) {
	assert.Error(t, OneOf[int]().Validate(0))
	assert.Error(t, OneOf(0).Validate(1))
	assert.Error(t, OneOf(0, 1, 2).Validate(3))
}
