package ini

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/zpatrick/cfg"
)

type Provider struct {
	f *ini.File
}

func New(path string) (*Provider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	return &Provider{f: f}, nil
}

func (p *Provider) String(section, key string) cfg.Provider[string] {
	return Provide(p, section, key, func(k *ini.Key) (string, error) { return k.String(), nil })
}

func (p *Provider) Float64(section, key string) cfg.Provider[float64] {
	return Provide(p, section, key, func(k *ini.Key) (float64, error) { return k.Float64() })
}

func (p *Provider) Int(section, key string) cfg.Provider[int] {
	return Provide(p, section, key, func(k *ini.Key) (int, error) { return k.Int() })
}

func (p *Provider) Int64(section, key string) cfg.Provider[int64] {
	return Provide(p, section, key, func(k *ini.Key) (int64, error) { return k.Int64() })
}

func (p *Provider) Uint64(section, key string) cfg.Provider[uint64] {
	return Provide(p, section, key, func(k *ini.Key) (uint64, error) { return k.Uint64() })
}

func (p *Provider) Bool(section, key string) cfg.Provider[bool] {
	return Provide(p, section, key, func(k *ini.Key) (bool, error) { return k.Bool() })
}

func (p *Provider) Duration(section, key string) cfg.Provider[time.Duration] {
	return Provide(p, section, key, func(k *ini.Key) (time.Duration, error) { return k.Duration() })
}

func Provide[T any](p *Provider, section string, key string, convert func(k *ini.Key) (T, error)) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		s, err := p.f.GetSection(section)
		if err != nil {
			if isSectionDoesNotExistErr(err, section) {
				return out, cfg.NoValueProvidedError
			}

			return out, err
		}

		k, err := s.GetKey(key)
		if err != nil {
			if isKeyDoesNotExistErr(err, key) {
				return out, cfg.NoValueProvidedError
			}

			return out, err
		}

		return convert(k)
	})
}

func isSectionDoesNotExistErr(err error, section string) bool {
	return strings.Contains(err.Error(), fmt.Sprintf("section %q does not exist", section))
}

func isKeyDoesNotExistErr(err error, key string) bool {
	return strings.Contains(err.Error(), fmt.Sprintf("key %q not exists", key))
}