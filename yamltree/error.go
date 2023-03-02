package yamltree

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ErrorWrongNodeType struct {
	want yaml.Kind
	got  yaml.Kind
}

func (err ErrorWrongNodeType) Error() string {
	return fmt.Sprintf("incorrect yaml node kind want %d got %d", err.want, err.got)
}
