package yamlconv

import (
	"gopkg.in/yaml.v3"
)

// MakeMap makes a map out of a yaml.MappingNode using the given functions
func MakeMap[K comparable, V any](node *yaml.Node, parseKey func(n *yaml.Node) (K, error), parseValue func(n *yaml.Node) (V, error)) (map[K]V, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 1 {
			return MakeMap(node.Content[0], parseKey, parseValue)
		}
	case yaml.MappingNode:
		m := make(map[K]V)
		for i := 0; i+1 < len(node.Content); i += 2 {
			k, err := parseKey(node.Content[i])
			if err != nil {
				return nil, err
			}
			v, err := parseValue(node.Content[i+1])
			if err != nil {
				return nil, err
			}
			m[k] = v
		}
		return m, nil
	}
	return nil, ErrorWrongNodeType{want: yaml.MappingNode, got: node.Kind}
}
