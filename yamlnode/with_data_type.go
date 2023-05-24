package yamlnode

import (
	"gopkg.in/yaml.v3"
)

type WithDataType[T any] struct {
	Data T
	Node *yaml.Node
}

func (wt *WithDataType[T]) UnmarshalYAML(value *yaml.Node) error {
	wt.Node = value
	return value.Decode(&wt.Data)
}

func (wt WithDataType[T]) MarshalYAML() (any, error) {
	return wt.Data, nil
}
