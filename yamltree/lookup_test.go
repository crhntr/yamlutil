package yamltree_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/crhntr/yamlutil/yamltree"
)

func Test_LookupKey(t *testing.T) {
	tests := []struct {
		Name        string
		InputYAML   string
		InputKey    string
		ExpectFound bool
		ExpectValue any
	}{
		{"not a hash", `[]`, "x", false, nil},
		{"empty hash", `{}`, "x", false, nil},
		{"found key", `{"banana": 32, "orange": 1}`, "banana", true, 32},
		{"not found key", `{"orange": 1}`, "banana", false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var node yaml.Node
			err := yaml.Unmarshal([]byte(tt.InputYAML), &node)
			require.NoError(t, err)

			require.NotPanics(t, func() {
				v, found := yamltree.LookupKey(&node, tt.InputKey)
				assert.Equal(t, found, tt.ExpectFound)

				if tt.ExpectValue != nil {
					valBuf, err := yaml.Marshal(v)
					require.NoError(t, err)
					expBuf, err := yaml.Marshal(tt.ExpectValue)
					require.NoError(t, err)
					assert.Equal(t, string(valBuf), string(expBuf))
				}
			})
		})
	}
}
