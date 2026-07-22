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
	if a == nil || b == nil {
		if a != b {
			t.Errorf("nil node mismatch at %s: a is nil: %t, b is nil: %t", p, a == nil, b == nil)
		}
		return
	}
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
		if !assert.Equalf(t, len(a.Content), len(b.Content), "mismatched document content lengths at %s", p) {
			return
		}
		for i := range a.Content {
			assertEqual(t, a.Content[i], b.Content[i], p)
		}
	case yaml.MappingNode:
		aKeys := slices.Collect(yamlnode.ValuesStrings(slices.Collect(yamlnode.Keys(a))))
		bKeys := slices.Collect(yamlnode.ValuesStrings(slices.Collect(yamlnode.Keys(b))))
		if !assert.Equal(t, aKeys, bKeys, p) {
			return
		}
		aValues := mappingValues(a)
		bValues := mappingValues(b)
		for i, key := range aKeys {
			assertEqual(t, aValues[i], bValues[i], fmt.Sprintf("%s.[%q]", p, key))
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
		t.Errorf("un-supported node kind %d at %s", a.Kind, p)
	}
}

func mappingValues(node *yaml.Node) []*yaml.Node {
	var values []*yaml.Node
	for _, value := range yamlnode.KeyValue(node) {
		values = append(values, value)
	}
	return values
}
