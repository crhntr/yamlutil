package yamltest

import (
	"fmt"
	"testing"

	"github.com/crhntr/yamlutil/yamltree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

func AssertEqual(t *testing.T, a, b *yaml.Node) {
	t.Helper()
	assertEqual(t, a, b, "")
}

func assertEqual(t *testing.T, a, b *yaml.Node, p string) {
	assert.Equalf(t, a.Kind, b.Kind, "kind mismatch at %s", p)
	if t.Failed() {
		return
	}
	switch a.Kind {
	case yaml.DocumentNode:
		assertEqual(t, a.Content[0], b.Content[0], "")
	case yaml.MappingNode:
		aKeys := yamltree.ValuesStrings(yamltree.Keys(a))
		bKeys := yamltree.ValuesStrings(yamltree.Keys(b))
		assert.Equal(t, aKeys, bKeys, p)
		keys := append(aKeys[:len(aKeys):len(aKeys)], bKeys...)
		slices.Sort(keys)
		for _, key := range slices.Compact(keys) {
			aValue, aFound := yamltree.LookupKey(a, key)
			assert.Truef(t, aFound, "a missing key %q in %s", key, p)
			bValue, bFound := yamltree.LookupKey(b, key)
			assert.Truef(t, bFound, "b missing key %q in %s", key, p)
			if aFound && bFound {
				assertEqual(t, aValue, bValue, fmt.Sprintf("%s.[%q]", p, key))
			}
		}
	case yaml.SequenceNode:
		require.Equal(t, len(a.Content), len(b.Content), "mismatched sequence lengths at %s", p)
		for i := 0; i < len(a.Content); i++ {
			assertEqual(t, a.Content[i], b.Content[i], fmt.Sprintf("%s[%d]", p, i))
		}
	case yaml.ScalarNode:
		assert.Equalf(t, a.Value, b.Value, "mismatched value at %s", p)
		assert.Equalf(t, a.Tag, b.Tag, "mismatched tag at %s", p)
	default:
		panic(fmt.Sprintf("un-supported node type %d at %s", a.Kind, p))
	}
}
