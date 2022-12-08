package cfg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/zpatrick/cfg"
)

type INIFileProvider struct {
	f *ini.File
}

func INIFile(path string) (*INIFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	return &INIFileProvider{f: f}, nil
}

func (i *INIFileProvider) String(section, key string) cfg.Provider[string] {
	return Provide(i, section, key, func(k *ini.Key) (string, error) { return k.String(), nil })
}

func (i *INIFileProvider) Float64(section, key string) cfg.Provider[float64] {
	return Provide(i, section, key, func(k *ini.Key) (float64, error) { return k.Float64() })
}

func (i *INIFileProvider) Int(section, key string) cfg.Provider[int] {
	return Provide(i, section, key, func(k *ini.Key) (int, error) { return k.Int() })
}

func (i *INIFileProvider) Int64(section, key string) cfg.Provider[int64] {
	return Provide(i, section, key, func(k *ini.Key) (int64, error) { return k.Int64() })
}

func (i *INIFileProvider) Uint64(section, key string) cfg.Provider[uint64] {
	return Provide(i, section, key, func(k *ini.Key) (uint64, error) { return k.Uint64() })
}

func (i *INIFileProvider) Bool(section, key string) cfg.Provider[bool] {
	return Provide(i, section, key, func(k *ini.Key) (bool, error) { return k.Bool() })
}

func (i *INIFileProvider) Duration(section, key string) cfg.Provider[time.Duration] {
	return Provide(i, section, key, func(k *ini.Key) (time.Duration, error) { return k.Duration() })
}

func Provide[T any](i *INIFileProvider, section string, key string, convert func(k *ini.Key) (T, error)) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		s, err := i.f.GetSection(section)
		if err != nil {
			if isINISectionDoesNotExistErr(err, section) {
				return out, cfg.NoValueProvidedError
			}

			return out, err
		}

		k, err := s.GetKey(key)
		if err != nil {
			if isINIKeyDoesNotExistErr(err, key) {
				return out, cfg.NoValueProvidedError
			}

			return out, err
		}

		return convert(k)
	})
}

func isINISectionDoesNotExistErr(err error, section string) bool {
	return strings.Contains(err.Error(), fmt.Sprintf("section %q does not exist", section))
}

func isINIKeyDoesNotExistErr(err error, key string) bool {
	return strings.Contains(err.Error(), fmt.Sprintf("key %q not exists", key))
}
