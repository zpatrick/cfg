package cfg

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-ini/ini"
	"sigs.k8s.io/yaml"
)

// A Node ...
type Node struct {
	Name     string
	Value    any
	Children map[string]*Node
}

// A FileContentParser converts raw content into a configuration tree...
type FileContentParser interface {
	Parse(content []byte) (*Node, error)
}

// The FileContentParserFunc is an adapter type which allows ordinary functions to be used as FileContentParsers.
type FileContentParserFunc func(content []byte) (*Node, error)

// Parse calls f(content).
func (f FileContentParserFunc) Parse(content []byte) (*Node, error) {
	return f(content)
}

func ParseJSON() FileContentParser {
	return FileContentParserFunc(func(content []byte) (*Node, error) {
		var unmarshaled map[string]any
		if err := json.Unmarshal(content, &unmarshaled); err != nil {
			return nil, err
		}

		root := &Node{
			Name:     "root",
			Children: createChildNodes(unmarshaled),
		}

		return root, nil
	})
}

func createChildNodes(values map[string]interface{}) map[string]*Node {
	nodes := make(map[string]*Node, len(values))
	for k, v := range values {
		child := &Node{Name: k, Value: v}

		switch v := v.(type) {
		case map[string]interface{}:
			child.Children = createChildNodes(v)
		}

		nodes[k] = child
	}

	return nodes
}

func ParseYAML() FileContentParser {
	return FileContentParserFunc(func(content []byte) (*Node, error) {
		asJSON, err := yaml.YAMLToJSON(content)
		if err != nil {
			return nil, err
		}

		return ParseJSON().Parse(asJSON)
	})
}

// Note that the underlying parser will parse all types as strings.
func ParseINI() FileContentParser {
	return FileContentParserFunc(func(content []byte) (*Node, error) {
		f, err := ini.Load(content)
		if err != nil {
			return nil, err
		}

		sectionNodes := map[string]*Node{}
		for _, s := range f.Sections() {
			childNodes := map[string]*Node{}
			for _, k := range s.Keys() {
				childNode := &Node{
					Name:  k.Name(),
					Value: k.Value(),
				}

				childNodes[k.Name()] = childNode
			}

			sectionNode := &Node{
				Name:     s.Name(),
				Children: childNodes,
			}

			sectionNodes[s.Name()] = sectionNode
		}

		root := &Node{
			Name:     "root",
			Children: sectionNodes,
		}

		return root, nil
	})
}

type FileProvider struct {
	root *Node
}

func File(parser FileContentParser, path string) (*FileProvider, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	root, err := parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return &FileProvider{root: root}, nil
}

func (f *FileProvider) Int(section string, keys ...string) Provider[int] {
	return FileProvide[int](f, section, keys...)
}

func (f *FileProvider) Int64(section string, keys ...string) Provider[int64] {
	return FileProvide[int64](f, section, keys...)
}

func (f *FileProvider) Float64(section string, keys ...string) Provider[float64] {
	return FileProvide[float64](f, section, keys...)
}

func (f *FileProvider) String(section string, keys ...string) Provider[string] {
	return FileProvide[string](f, section, keys...)
}

func (f *FileProvider) Bool(section string, keys ...string) Provider[bool] {
	return FileProvide[bool](f, section, keys...)
}

// TOOD: uints,  float, etc.

// FileProvide
// Generics are not supported on methods (yet), hence why this is a package function instead of a method.
func FileProvide[T any](f *FileProvider, section string, keys ...string) Provider[T] {
	return ProviderFunc[T](func(ctx context.Context) (out T, err error) {
		n, ok := f.root.Children[section]
		if !ok {
			return out, NoValueProvidedError
		}

		for _, key := range keys {
			n, ok = n.Children[key]
			if !ok {
				return out, NoValueProvidedError
			}
		}

		out, ok = n.Value.(T)
		if !ok {
			return out, NewUnexpectedTypeError(out, n.Value)
		}

		return out, nil
	})
}
