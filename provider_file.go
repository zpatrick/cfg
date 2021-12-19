package cfg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type FileFormat string

const (
	FormatINI  FileFormat = "ini"
	FormatYAML FileFormat = "yaml"
	FormatJSON FileFormat = "json"
)

type fileProvider struct {
	root *node
}

func (f *fileProvider) Provide(section string, keys ...string) Provider {
	return ProviderFunc(func(ctx context.Context) ([]byte, error) {
		n, ok := f.root.children[section]
		if !ok {
			return nil, NoValueProvidedError
		}

		for _, key := range keys {
			n, ok = n.children[key]
			if !ok {
				return nil, NoValueProvidedError
			}
		}

		return Encode(n.value), nil
	})
}

func File(format FileFormat, path string) (*fileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var parseFunc func(data []byte) (*node, error)
	switch format {
	case FormatINI:
	case FormatYAML:
	case FormatJSON:
		parseFunc = parseJSON
	default:
		return nil, fmt.Errorf("unrecognized file format: %v", format)
	}

	root, err := parseFunc(data)
	if err != nil {
		return nil, err
	}

	return &fileProvider{root: root}, nil
}

type node struct {
	name     string
	value    interface{}
	children map[string]*node
}

func (n *node) IsLeaf() bool {
	return len(n.children) == 0
}

func parseJSON(data []byte) (*node, error) {
	unmarshaled := map[string]interface{}{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		return nil, err
	}

	root := &node{
		name:     "root",
		children: createChildNodes(unmarshaled),
	}

	return root, nil
}

func createChildNodes(values map[string]interface{}) map[string]*node {
	nodes := make(map[string]*node, len(values))
	for k, v := range values {
		child := &node{name: k}

		switch v := v.(type) {
		case map[string]interface{}:
			child.children = createChildNodes(v)
		default:
			child.value = v
		}

		nodes[k] = child
	}

	return nodes
}
