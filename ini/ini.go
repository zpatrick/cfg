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

// A Provider loads configuration values from a file in ini format.
type Provider struct {
	f *ini.File
}

// Load reads and parses the ini file at the given path.
func Load(path string) (*Provider, error) {
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

// New returns a Provider which loads values from f.
func New(f *ini.File) (*Provider, error) {
	return &Provider{f: f}, nil
}

// String returns a string configuration value at the given section and key.
func (p *Provider) String(section, key string) cfg.Provider[string] {
	return Provide(p, section, key, func(k *ini.Key) (string, error) { return k.String(), nil })
}

// Float64 returns a float64 configuration value at the given section and key.
func (p *Provider) Float64(section, key string) cfg.Provider[float64] {
	return Provide(p, section, key, func(k *ini.Key) (float64, error) { return k.Float64() })
}

// Int returns an int configuration value at the given section and key.
func (p *Provider) Int(section, key string) cfg.Provider[int] {
	return Provide(p, section, key, func(k *ini.Key) (int, error) { return k.Int() })
}

// Int64 returns an int64 configuration value at the given section and key.
func (p *Provider) Int64(section, key string) cfg.Provider[int64] {
	return Provide(p, section, key, func(k *ini.Key) (int64, error) { return k.Int64() })
}

// UInt64 returns an uint64 configuration value at the given section and key.
func (p *Provider) Uint64(section, key string) cfg.Provider[uint64] {
	return Provide(p, section, key, func(k *ini.Key) (uint64, error) { return k.Uint64() })
}

// Bool returns a bool configuration value at the given section and key.
func (p *Provider) Bool(section, key string) cfg.Provider[bool] {
	return Provide(p, section, key, func(k *ini.Key) (bool, error) { return k.Bool() })
}

// Duration returns a duration configuration value at the given section and key.
func (p *Provider) Duration(section, key string) cfg.Provider[time.Duration] {
	return Provide(p, section, key, func(k *ini.Key) (time.Duration, error) { return k.Duration() })
}

// Provide loads a *ini.Key from the given section and key from p.
// The convert paramter is then called to convert the *ini.Key into a T.
// If the key does not exist, a cfg.NoValueProvidedError error will be returned.
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
