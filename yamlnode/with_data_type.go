package yamlnode

import (
	"gopkg.in/yaml.v3"
)

// WithDataType is a wrapper for a yaml.Node that also stores the decoded value.
// It is useful when you want to reflect back on the original yaml.Node after parsing.
// This may be helpful for providing richer error messages for configuration problems.
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
