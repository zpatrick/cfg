package cfg

import (
	"os"

	"github.com/BurntSushi/toml"
)

type TOMLFileProvider struct {
	mapProvider
}

func TOMLFile(path string) (*TOMLFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := toml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &TOMLFileProvider{mapProvider(m)}, nil
}
