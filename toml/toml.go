package toml

import (
	"context"
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal"
)

type Provider struct {
	tree *toml.Tree
}

func New(path string) (*Provider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t, err := toml.LoadBytes(data)
	if err != nil {
		return nil, err
	}

	return &Provider{tree: t}, nil
}

func (p *Provider) String(keys ...string) cfg.Provider[string] {
	return Provide[string](p, keys...)
}

func (p *Provider) Int(keys ...string) cfg.Provider[int] {
	return Provide[int](p, keys...)
}

func (p *Provider) Int64(keys ...string) cfg.Provider[int64] {
	return Provide[int64](p, keys...)
}

func (p *Provider) Float64(keys ...string) cfg.Provider[float64] {
	return Provide[float64](p, keys...)
}

func (p *Provider) Bool(keys ...string) cfg.Provider[bool] {
	return Provide[bool](p, keys...)
}

func (p *Provider) Duration(keys ...string) cfg.Provider[time.Duration] {
	return cfg.ProviderFunc[time.Duration](func(ctx context.Context) (time.Duration, error) {
		val, err := Provide[string](p, keys...).Provide(ctx)
		if err != nil {
			return 0, err
		}

		return time.ParseDuration(val)
	})
}

func Provide[T any](p *Provider, keys ...string) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		if !p.tree.HasPath(keys) {
			return out, cfg.NoValueProvidedError
		}

		val := p.tree.GetPath(keys)
		v, ok := val.(T)
		if !ok {
			return out, internal.NewUnexpectedTypeError(out, val)
		}

		return v, nil
	})
}
