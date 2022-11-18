package yaml

import (
	"os"

	"github.com/go-yaml/yaml"
	"github.com/zpatrick/cfg/providers/generic"
)

type YAMLFileProvider struct {
	generic.Provider
}

func NewFile(path string) (*YAMLFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &YAMLFileProvider{generic.Provider(m)}, nil
}
