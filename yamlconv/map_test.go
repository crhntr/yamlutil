package yamlconv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/crhntr/yamlutil/yamlconv"
)

func TestMakeMap(t *testing.T) {
	parseString := func(n *yaml.Node) (string, error) {
		var s string
		err := n.Decode(&s)
		return s, err
	}
	parseInt := func(n *yaml.Node) (int, error) {
		var i int
		err := n.Decode(&i)
		return i, err
	}

	for _, tt := range []struct {
		Name      string
		InputYAML string
		Expect    map[string]int
		ExpectErr bool
	}{
		{Name: "mapping", InputYAML: `{a: 1, b: 2}`, Expect: map[string]int{"a": 1, "b": 2}},
		{Name: "empty mapping", InputYAML: `{}`, Expect: map[string]int{}},
		{Name: "duplicate keys", InputYAML: "a: 1\na: 2\n", ExpectErr: true},
		{Name: "sequence", InputYAML: `[1, 2]`, ExpectErr: true},
		{Name: "scalar", InputYAML: `5`, ExpectErr: true},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			var node yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(tt.InputYAML), &node))

			m, err := yamlconv.MakeMap(&node, parseString, parseInt)

			if tt.ExpectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.Expect, m)
		})
	}

	t.Run("document with multiple children", func(t *testing.T) {
		node := &yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{Kind: yaml.MappingNode},
				{Kind: yaml.MappingNode},
			},
		}
		_, err := yamlconv.MakeMap(node, parseString, parseInt)
		require.Error(t, err)
		assert.ErrorContains(t, err, "document")
	})
}
