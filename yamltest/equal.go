package yamltest

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"

	"github.com/crhntr/yamlutil/yamlnode"
)

// AssertEqual asserts that two yaml.Node trees are equal.
// It follows aliases.
func AssertEqual(t *testing.T, a, b *yaml.Node) {
	t.Helper()
	assertEqual(t, a, b, "")
}

func assertEqual(t require.TestingT, a, b *yaml.Node, p string) {
	if a.Kind == yaml.AliasNode {
		assertEqual(t, a.Alias, b, p)
		return
	}
	if b.Kind == yaml.AliasNode {
		assertEqual(t, a, b.Alias, p)
		return
	}

	require.Equalf(t, a.Kind, b.Kind, "kind mismatch at %s", p)
	switch a.Kind {
	case yaml.DocumentNode:
		assertEqual(t, a.Content[0], b.Content[0], "")
	case yaml.MappingNode:
		aKeys := slices.Collect(yamlnode.ValuesStrings(slices.Collect(yamlnode.Keys(a))))
		bKeys := slices.Collect(yamlnode.ValuesStrings(slices.Collect(yamlnode.Keys(b))))
		assert.Equal(t, aKeys, bKeys, p)
		keys := append(aKeys[:len(aKeys):len(aKeys)], bKeys...)
		slices.Sort(keys)
		for _, key := range slices.Compact(keys) {
			aValue, aFound := yamlnode.LookupKey(a, key)
			assert.Truef(t, aFound, "a missing key %q in %s", key, p)
			bValue, bFound := yamlnode.LookupKey(b, key)
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
	case 0:
		assert.Zero(t, *a)
		assert.Zero(t, *b)
	default:
		panic(fmt.Sprintf("un-supported node type %d at %s", a.Kind, p))
	}
}
