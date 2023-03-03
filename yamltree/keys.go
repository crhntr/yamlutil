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

func Keys(node *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	switch node.Kind {
	case yaml.DocumentNode:
		return Keys(node.Content[0])
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			result = append(result, node.Content[i])
		}
	}
	return result
}

func ValuesStrings(nodes []*yaml.Node) []string {
	result := make([]string, 0, len(nodes))
	for _, x := range nodes {
		result = append(result, x.Value)
	}
	return result
}
