package yamlnode

import "gopkg.in/yaml.v3"

// LookupKey does a simple lookup of a string key in a yaml.MappingNode
// it does not support recursive lookup, or any other fancy features.
// Consider using Walk or github.com/mikefarah/yq/v3/pkg/yqlib for more sophisticated queries.
func LookupKey(node *yaml.Node, key string) (*yaml.Node, bool) {
	if node == nil {
		return nil, false
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, c := range node.Content {
			found, isFound := LookupKey(c, key)
			if isFound {
				return found, isFound
			}
		}
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				return node.Content[i+1], true
			}
		}
	case yaml.ScalarNode, yaml.AliasNode, yaml.SequenceNode:
	}
	return nil, false
}

// Keys returns the keys of a yaml.MappingNode.
func Keys(node *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	switch node.Kind {
	case yaml.DocumentNode:
		return Keys(node.Content[0])
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			result = append(result, node.Content[i])
		}
	case yaml.ScalarNode, yaml.AliasNode, yaml.SequenceNode:
	}
	return result
}

// ValuesStrings returns the values of a list of nodes.
func ValuesStrings(nodes []*yaml.Node) []string {
	result := make([]string, 0, len(nodes))
	for _, x := range nodes {
		result = append(result, x.Value)
	}
	return result
}
