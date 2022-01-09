package cfg

import (
	"os"

	"github.com/go-yaml/yaml"
)

type JSONFileProvider struct {
	mapProvider
}

func JSONFile(path string) (*JSONFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m mapProvider
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &JSONFileProvider{m}, nil
}
