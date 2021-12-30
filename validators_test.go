package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetween(t *testing.T) {
	assert.NoError(t, Between(0, 2).Validate(1))
	assert.NoError(t, Between(-1, 1).Validate(0))
	assert.NoError(t, Between(8.5, 8.6).Validate(8.55))
}

func TestBetweenError(t *testing.T) {
	assert.Error(t, Between(0, 2).Validate(0))
	assert.Error(t, Between(0, 2).Validate(2))
	assert.Error(t, Between(0, 2).Validate(-1))
	assert.Error(t, Between(0, 2).Validate(3))
	assert.Error(t, Between(8.5, 8.6).Validate(8.4))

	Or(OneOf(22, 90), Between(8000, 9000))
}

func TestOneOf(t *testing.T) {
	assert.NoError(t, OneOf(0).Validate(0))
	assert.NoError(t, OneOf(0, 1).Validate(1))
	assert.NoError(t, OneOf(0, 1, 2).Validate(1))
	assert.NoError(t, OneOf(5, 5, 5).Validate(5))
}

func TestOneOfError(t *testing.T) {
	assert.Error(t, OneOf[int]().Validate(0))
	assert.Error(t, OneOf(0).Validate(1))
	assert.Error(t, OneOf(0, 1, 2).Validate(3))
}
