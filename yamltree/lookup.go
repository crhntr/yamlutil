package yamltree

import "gopkg.in/yaml.v3"

func LookupKey(node *yaml.Node, key string) (*yaml.Node, bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 1 {
			return LookupKey(node.Content[0], key)
		}
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				return node.Content[i+1], true
			}
		}
	}
	return nil, false
}
