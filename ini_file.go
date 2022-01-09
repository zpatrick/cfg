package cfg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
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

// TODO: int64, float
func (i *INIFileProvider) String(section, key string) Provider[string] {
	return INIFileProvide(i, section, key, func(k *ini.Key) (string, error) { return k.String(), nil })
}

func (i *INIFileProvider) Int(section, key string) Provider[int] {
	return INIFileProvide(i, section, key, func(k *ini.Key) (int, error) { return k.Int() })
}

func (i *INIFileProvider) Bool(section, key string) Provider[bool] {
	return INIFileProvide(i, section, key, func(k *ini.Key) (bool, error) { return k.Bool() })
}

func (i *INIFileProvider) Duration(section, key string) Provider[time.Duration] {
	return INIFileProvide(i, section, key, func(k *ini.Key) (time.Duration, error) { return k.Duration() })
}

func INIFileProvide[T any](i *INIFileProvider, section string, key string, convert func(k *ini.Key) (T, error)) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		s, err := i.f.GetSection(section)
		if err != nil {
			if isINISectionDoesNotExistErr(err, section) {
				return out, NoValueProvidedError
			}

			return out, err
		}

		k, err := s.GetKey(key)
		if err != nil {
			if isINIKeyDoesNotExistErr(err, key) {
				return out, NoValueProvidedError
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
