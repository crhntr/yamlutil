package yamlnode

import (
	"gopkg.in/yaml.v3"
)

// Walk walks a yaml.Node tree, calling fn on each node.
// If fn returns an error, Walk returns that error and stops walking.
// Aliases are followed and may cause an infinite loop.
func Walk(node *yaml.Node, fn func(node *yaml.Node) error, kinds yaml.Kind) error {
	if node.Kind&kinds != 0 {
		err := fn(node)
		if err != nil {
			return err
		}
	}
	switch node.Kind {
	case yaml.DocumentNode, yaml.MappingNode, yaml.SequenceNode:
		for _, child := range node.Content {
			err := Walk(child, fn, kinds)
			if err != nil {
				return err
			}
		}
	case yaml.AliasNode:
		err := Walk(node.Alias, fn, kinds)
		if err != nil {
			return err
		}
	case yaml.ScalarNode:
	}
	return nil
}
