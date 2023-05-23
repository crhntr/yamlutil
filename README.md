# yamlutil [![Go Reference](https://pkg.go.dev/badge/github.com/crhntr/yamlutil.svg)](https://pkg.go.dev/github.com/crhntr/yamlutil)

This repository has utilities for working with the Node type in [YAML V3](https://pkg.go.dev/gopkg.in/yaml.v3).

## The yamlconv package

- It has a function to convert a `*yaml.Node` to JSON.
- It also contains a function to convert a `*yaml.Node` with `yaml.Kind` `yaml.DocumentNode` or  `yaml.MappingNode` to a Go map. 

## The yamlnode package

Contains a Walk function to traverse a YAML document.
It also contains, a method to lookup the value given a (string) key in a `yaml.DocumentNode` or  `yaml.MappingNode`.

You may want to use [yqlib](https://pkg.go.dev/github.com/mikefarah/yq/v4@v4.31.2/pkg/yqlib) instead.
