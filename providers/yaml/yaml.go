package yaml

import (
	"context"
	"os"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/providers/generic"
)

type Provider struct {
	provider generic.Provider
}

func New(path string) (*Provider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &Provider{generic.Provider(m)}, nil
}

func (p *Provider) String(section string, keys ...string) cfg.Provider[string] {
	return p.provider.String(section, keys...)
}

func (p *Provider) Float64(section string, keys ...string) cfg.Provider[float64] {
	return p.provider.Float64(section, keys...)
}

func (p *Provider) Int(section string, keys ...string) cfg.Provider[int] {
	return p.provider.Int(section, keys...)
}

func (p *Provider) Bool(section string, keys ...string) cfg.Provider[bool] {
	return p.provider.Bool(section, keys...)
}

func (p *Provider) Duration(section string, keys ...string) cfg.Provider[time.Duration] {
	return cfg.ProviderFunc[time.Duration](func(ctx context.Context) (out time.Duration, err error) {
		val, err := Provide[string](p, section, keys...).Provide(ctx)
		if err != nil {
			return out, err
		}

		return time.ParseDuration(val)
	})
}

func Provide[T any](p *Provider, section string, keys ...string) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		val, err := p.provider.Get(section, keys...)
		if err != nil {
			return out, cfg.NoValueProvidedError
		}

		v, ok := val.(T)
		if !ok {
			return out, generic.NewUnexpectedTypeError(out, val)
		}

		return v, nil
	})
}
