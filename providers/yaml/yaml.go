package yaml

import (
	"context"
	"os"

	"github.com/zpatrick/cfg"
	"gopkg.in/yaml.v3"
)

// TODO: comments
type Provider struct {
	root *yaml.Node
}

func New(path string) (*Provider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return &Provider{root: &root}, nil
}

func (p *Provider) Bool(keys ...string) cfg.Provider[bool] {
	return Provide[bool](p, keys...)
}

func Provide[T any](p *Provider, keys ...string) cfg.Provider[T] {
	return cfg.ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		r := p.root

		print(r)
		return out, nil
	})
}

// unmarhsla; https://github.com/go-yaml/yaml/blob/v3.0.1/decode.go#L484
