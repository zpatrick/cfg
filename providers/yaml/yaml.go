package yaml

import (
	"os"

	"github.com/zpatrick/cfg/providers/generic"
	"gopkg.in/yaml.v3"
)

// TODO: comments
type Provider struct {
	generic.Provider
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
