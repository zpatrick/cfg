package cfg

import (
	"encoding/json"
	"os"
)

type JSONFileProvider struct {
	mapProvider
}

func JSONFile(path string) (*JSONFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &JSONFileProvider{mapProvider(m)}, nil
}
