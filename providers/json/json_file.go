package cfg

import (
	"encoding/json"
	"os"

	"github.com/zpatrick/cfg/providers/generic"
)

type JSONFileProvider struct {
	generic.Provider
}

func New(path string) (*JSONFileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &JSONFileProvider{generic.Provider(m)}, nil
}
