package yaml

import (
	"os"
	"time"

	"github.com/zpatrick/cfg"
	"github.com/zpatrick/cfg/internal"
	"gopkg.in/yaml.v3"
)

type Provider struct {
	m internal.MapProvider
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

	return &Provider{internal.MapProvider(m)}, nil
}

func (p Provider) String(section string, keys ...string) cfg.Provider[string] {
	return p.m.String(section, keys...)
}

func (p Provider) Int(section string, keys ...string) cfg.Provider[int] {
	return p.m.Int(section, keys...)
}

func (p Provider) Bool(section string, keys ...string) cfg.Provider[bool] {
	return p.m.Bool(section, keys...)
}

func (p Provider) Duration(section string, keys ...string) cfg.Provider[time.Duration] {
	return p.m.Duration(section, keys...)
}

func (p Provider) Float64(section string, keys ...string) cfg.Provider[float64] {
	return p.m.Float64(section, keys...)
}
