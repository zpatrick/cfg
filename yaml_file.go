package cfg

import (
	"os"

	"github.com/go-yaml/yaml"
)

type YAMLFileProvider struct {
	mapProvider
}

func YAMLFile(path string) (*YAMLFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m mapProvider
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &YAMLFileProvider{m}, nil
}
