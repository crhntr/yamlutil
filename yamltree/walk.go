package yamltree

import (
	"gopkg.in/yaml.v3"
)

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
