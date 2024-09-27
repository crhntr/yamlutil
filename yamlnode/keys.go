package yamlnode

import (
	"iter"

	"gopkg.in/yaml.v3"
)

// LookupKey does a simple lookup of a string key in a yaml.MappingNode
// it does not support recursive lookup, or any other fancy features.
// Consider using Walk or github.com/mikefarah/yq/v3/pkg/yqlib for more sophisticated queries.
func LookupKey(node *yaml.Node, key string) (*yaml.Node, bool) {
	return LookupValueFunc(node, func(n *yaml.Node) bool {
		return n.Value == key
	})
}

// LookupValueFunc does a simple lookup of a node in a yaml.MappingNode
// It does not support recursive lookup, or any other fancy features.
// It uses KeyValue and will iterate yaml.DocumentNode and yaml.AliasNode.
func LookupValueFunc(node *yaml.Node, checkKey func(n *yaml.Node) bool) (*yaml.Node, bool) {
	for k, v := range KeyValue(node) {
		if checkKey(k) {
			return v, true
		}
	}
	return nil, false
}

// Keys returns the keys of a yaml.MappingNode.
// It uses KeyValue and will iterate yaml.DocumentNode and yaml.AliasNode.
func Keys(node *yaml.Node) iter.Seq[*yaml.Node] {
	return func(yield func(*yaml.Node) bool) {
		for k, _ := range KeyValue(node) {
			if !yield(k) {
				return
			}
		}
	}
}

// ValuesStrings returns the Value fields from a yaml.Node of a list of nodes.
func ValuesStrings(nodes []*yaml.Node) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, n := range nodes {
			if !yield(n.Value) {
				return
			}
		}
	}
}

// KeyValue iterates over a yaml.MappingNode kinda like a map[*yaml.Node][*yaml.Node]
// It iterates multiple content elements in yaml.DocumentNode and will follow yaml.AliasNode.
func KeyValue(node *yaml.Node) iter.Seq2[*yaml.Node, *yaml.Node] {
	if node == nil {
		return func(func(*yaml.Node, *yaml.Node) bool) {}
	}
	switch node.Kind {
	case yaml.DocumentNode:
		return func(yield func(*yaml.Node, *yaml.Node) bool) {
			for _, child := range node.Content {
				for k, v := range KeyValue(child) {
					if !yield(k, v) {
						return
					}
				}
			}
		}
	case yaml.MappingNode:
		return func(yield func(*yaml.Node, *yaml.Node) bool) {
			for i := 0; i+1 < len(node.Content); i += 2 {
				if !yield(node.Content[i], node.Content[i+1]) {
					return
				}
			}
		}
	case yaml.AliasNode:
		return KeyValue(node.Alias)
	default:
		return func(func(*yaml.Node, *yaml.Node) bool) {}
	}
}
