package yamlconv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/crhntr/yamlutil/yamlconv"
)

func TestToJSON(t *testing.T) {
	t.Run("nil node", func(t *testing.T) {
		require.NotPanics(t, func() {
			_, _ = yamlconv.ToJSON(nil, nil)
		})
	})

	for _, tt := range []struct {
		Name,
		inYAML, outJSON string
	}{
		{
			Name:   "empty",
			inYAML: ` `, outJSON: "{}",
		},
		{
			Name: "list",
			inYAML: `---
- apple
- banana
`, outJSON: `["apple","banana"]`,
		},
		{
			Name: "null values",
			inYAML: `---
-
-
`, outJSON: `[null,null]`,
		},
		{
			Name:   "empty array",
			inYAML: `[]`, outJSON: `[]`,
		},
		{
			Name:   "empty curly map",
			inYAML: `{}`, outJSON: `{}`,
		},
		{
			Name:   "empty document",
			inYAML: `---`, outJSON: `{}`,
		},
		{
			Name: "maps",
			inYAML: `---
value: x
emptyMap: {}
mapWithValue: {x: y}
`, outJSON: `{"value":"x","emptyMap":{},"mapWithValue":{"x":"y"}}`,
		},
		{
			Name: "map with quoted keys",
			inYAML: `---
"key1": 0
'key2': 0
key3: 0
`, outJSON: `{"key1":0,"key2":0,"key3":0}`,
		},
		{
			Name: "some zeros",
			inYAML: `---
- 0
- 0.0
- 0.00000000
`, outJSON: `[0,0.0,0.0]`,
		},
		{
			Name: "mixed integers and floats",
			inYAML: `---
- 1
- 1.1
- 1.0
`, outJSON: `[1,1.1,1.0]`,
		},
		{
			Name: "some other type",
			inYAML: `---
didItHappen: true
`, outJSON: `{"didItHappen":true}`,
		},
		{
			Name: "alias",
			inYAML: `---
x: &val 4
y: *val
`, outJSON: `{"x":4,"y":4}`,
		},
		{
			Name: "embedded",
			inYAML: `---
x:
  y: 5
`, outJSON: `{"x":{"y":5}}`,
		},
		{
			Name: "single space indent",
			inYAML: `---
x:
 y:
  z: 7
`, outJSON: `{"x":{"y":{"z":7}}}`,
		},
		{
			Name:    "tagged binary",
			inYAML:  `picture: !!binary ` + somePNG,
			outJSON: `{"picture":"` + somePNG + `"}`,
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			var n yaml.Node
			err := yaml.Unmarshal([]byte(tt.inYAML), &n)
			require.NoError(t, err)
			out, err := yamlconv.ToJSON(nil, &n)
			assert.NoError(t, err)
			assert.Equal(t, tt.outJSON, string(out))
		})
	}
}

const (
	somePNG = `iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAAA1BMVEW10NBjBBbqAAAAH0lEQVRoge3BAQ0AAADCoPdPbQ43oAAAAAAAAAAAvg0hAAABmmDh1QAAAABJRU5ErkJggg==`
)
