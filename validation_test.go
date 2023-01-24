package cfg_test

import (
	"testing"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/testx/assert"
)

func TestBetween(t *testing.T) {
	assert.Error(t, cfg.Between(0, 2).Validate(-1))
	assert.NilError(t, cfg.Between(0, 2).Validate(0))
	assert.NilError(t, cfg.Between(0, 2).Validate(1))
	assert.NilError(t, cfg.Between(0, 2).Validate(2))
	assert.Error(t, cfg.Between(0, 2).Validate(3))
}

func TestOneOf(t *testing.T) {
	assert.Error(t, cfg.OneOf[int]().Validate(0))
	assert.Error(t, cfg.OneOf(1, 2, 3).Validate(0))
	assert.NilError(t, cfg.OneOf(1, 2, 3).Validate(1))
	assert.NilError(t, cfg.OneOf(1, 2, 3).Validate(2))
	assert.NilError(t, cfg.OneOf(1, 2, 3).Validate(3))
	assert.Error(t, cfg.OneOf(1, 2, 3).Validate(4))
}

func TestAnd(t *testing.T) {
	assert.NilError(t, cfg.And(
		cfg.Between(0, 100),
		cfg.Between(0, 50),
	).Validate(49))

	assert.Error(t, cfg.And(
		cfg.Between(0, 100),
		cfg.Between(0, 50),
	).Validate(51))
}

func TestOr(t *testing.T) {
	assert.NilError(t, cfg.Or(
		cfg.Between(0, 100),
		cfg.Between(0, 50),
	).Validate(51))

	assert.NilError(t, cfg.Or(
		cfg.Between(0, 50),
		cfg.Between(0, 100),
	).Validate(51))

	assert.Error(t, cfg.Or(
		cfg.Between(0, 100),
		cfg.Between(0, 50),
	).Validate(101))
}

func TestNot(t *testing.T) {
	assert.NilError(t, cfg.Not(cfg.OneOf(1, 2, 3)).Validate(0))
	assert.Error(t, cfg.Not(cfg.OneOf(1, 2, 3)).Validate(1))
}
