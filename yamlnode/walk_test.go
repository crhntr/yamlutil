package yamlnode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/crhntr/yamlutil/yamlnode"
)

func TestWalk(t *testing.T) {
	t.Run("nil node", func(t *testing.T) {
		require.NotPanics(t, func() {
			err := yamlnode.Walk(nil, func(*yaml.Node) error {
				t.Error("fn must not be called for a nil node")
				return nil
			}, yaml.ScalarNode)
			assert.NoError(t, err)
		})
	})

	t.Run("scalars in a mapping", func(t *testing.T) {
		var node yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte("a: 1\nb: 2\n"), &node))

		var values []string
		err := yamlnode.Walk(&node, func(n *yaml.Node) error {
			values = append(values, n.Value)
			return nil
		}, yaml.ScalarNode)

		require.NoError(t, err)
		assert.Equal(t, []string{"a", "1", "b", "2"}, values)
	})
}
