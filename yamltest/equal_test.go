package yamltest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestAssertEqual(t *testing.T) {
	const (
		embeddedMultiLine = `---
w:
  y:
    z: "banana"
x:
  y:
    z: "banana"
`
		anchorMultiLine = `---
w: &bananaYZ
  y:
    z: "banana"
x: *bananaYZ
`
		alternativeAnchor = `---
w:
  y: &bananaZ
    z: "banana"
x:
  y: *bananaZ
`
	)
	for _, tt := range []struct {
		Name       string
		A, B       string
		ShouldFail bool
	}{
		{Name: "empty", ShouldFail: false},
		{Name: "empty json", A: `{}`, B: `{}`, ShouldFail: false},
		{Name: "empty yaml", A: `---`, B: `---`, ShouldFail: false},
		{Name: "number mapping values", A: `x: 1`, B: `x: 1`, ShouldFail: false},
		{Name: "string mapping values", A: `x:"hello"`, B: `x:"hello"`, ShouldFail: false},
		{Name: "boolean true mapping values", A: `x:true`, B: `x:true`, ShouldFail: false},
		{Name: "boolean false mapping values", A: `x:false`, B: `x:false`, ShouldFail: false},
		{Name: "embedded inline", A: `x:{y:"z"}`, B: `x:{y:"z"}`, ShouldFail: false},
		{Name: "embedded multiline", A: embeddedMultiLine, B: embeddedMultiLine, ShouldFail: false},
		{Name: "different spacing inline", A: `x:    1       `, B: `x: 1`, ShouldFail: false},
		{Name: "empty sequences", A: `[]`, B: `[]`, ShouldFail: false},
		{Name: "one value in sequence", A: `x: [2]`, B: `x: [2]`, ShouldFail: false},
		{Name: "some values in sequence", A: `x: [1, 2, 4, 8]`, B: `x: [1,2,4,8]`, ShouldFail: false},
		{Name: "some values in sequence", A: `x: [1, 2, 4, 8]`, B: `x: [1,2,4,8]`, ShouldFail: false},
		{Name: "a is an alias", A: anchorMultiLine, B: embeddedMultiLine, ShouldFail: false},
		{Name: "b is an alias", A: embeddedMultiLine, B: anchorMultiLine, ShouldFail: false},
		{Name: "alternative aliases", A: alternativeAnchor, B: anchorMultiLine, ShouldFail: false},
		{Name: "alternative aliases swapped", A: anchorMultiLine, B: alternativeAnchor, ShouldFail: false},

		{Name: "equal duplicate keys", A: "a: 1\na: 2\n", B: "a: 1\na: 2\n", ShouldFail: false},

		{Name: "different spacing in quotes", A: `x: " banana"`, B: `x: "banana"`, ShouldFail: true},
		{Name: "duplicate keys with different later values", A: "a: 1\na: 2\n", B: "a: 1\na: 3\n", ShouldFail: true},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			var a, b yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(tt.A), &a))
			require.NoError(t, yaml.Unmarshal([]byte(tt.B), &b))
			mock := new(mockT)
			assertEqual(mock, &a, &b, "")
			if tt.ShouldFail {
				require.True(t, mock.Failed, "it should fail but got: %s", fmt.Sprintf(mock.format, mock.args))
			} else {
				require.False(t, mock.Failed, "it should succeed but got: %s", fmt.Sprintf(mock.format, mock.args))
			}
		})
	}

	t.Run("nil nodes are equal", func(t *testing.T) {
		mock := new(mockT)
		require.NotPanics(t, func() {
			assertEqual(mock, nil, nil, "")
		})
		require.False(t, mock.Failed)
	})

	t.Run("nil and non-nil nodes are not equal", func(t *testing.T) {
		mock := new(mockT)
		require.NotPanics(t, func() {
			assertEqual(mock, nil, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"}, "")
		})
		require.True(t, mock.Failed)
	})

	t.Run("alias with nil target does not panic", func(t *testing.T) {
		mock := new(mockT)
		require.NotPanics(t, func() {
			assertEqual(mock, &yaml.Node{Kind: yaml.AliasNode}, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"}, "")
		})
		require.True(t, mock.Failed)
	})

	t.Run("documents without children are equal", func(t *testing.T) {
		mock := new(mockT)
		require.NotPanics(t, func() {
			assertEqual(mock, &yaml.Node{Kind: yaml.DocumentNode}, &yaml.Node{Kind: yaml.DocumentNode}, "")
		})
		require.False(t, mock.Failed)
	})

	t.Run("documents with different child counts are not equal", func(t *testing.T) {
		mock := new(mockT)
		one := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"},
		}}
		require.NotPanics(t, func() {
			assertEqual(mock, &yaml.Node{Kind: yaml.DocumentNode}, one, "")
		})
		require.True(t, mock.Failed)
	})

	t.Run("documents differing after the first child are not equal", func(t *testing.T) {
		mock := new(mockT)
		a := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"},
			{Kind: yaml.ScalarNode, Tag: "!!int", Value: "2"},
		}}
		b := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"},
			{Kind: yaml.ScalarNode, Tag: "!!int", Value: "3"},
		}}
		assertEqual(mock, a, b, "")
		require.True(t, mock.Failed)
	})

	t.Run("unknown node kind fails instead of panicking", func(t *testing.T) {
		mock := new(mockT)
		require.NotPanics(t, func() {
			assertEqual(mock, &yaml.Node{Kind: yaml.Kind(99)}, &yaml.Node{Kind: yaml.Kind(99)}, "")
		})
		require.True(t, mock.Failed)
	})

	t.Run("exported", func(t *testing.T) {
		var a, b yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(`{}`), &a))
		require.NoError(t, yaml.Unmarshal([]byte(`{}`), &b))
		AssertEqual(t, &a, &b)
	})
}

type mockT struct {
	Failed bool
	format string
	args   []any
}

func (t *mockT) FailNow() {
	t.Failed = true
}

func (t *mockT) Errorf(format string, args ...any) {
	t.Failed = true
	t.format, t.args = format, args
}
