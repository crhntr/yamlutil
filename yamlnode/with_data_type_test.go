package yamlnode_test

import (
	"net/url"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crhntr/yamlutil/yamlnode"
)

func TestWithDataType(t *testing.T) {
	const in = `---
backend: "https://example.com:27017"
database: "test"
retry_count: 5
`

	type Configuration struct {
		yamlnode.WithDataType[struct {
			Database   string                        `yaml:"database"`
			Backend    yamlnode.WithDataType[Server] `yaml:"backend"`
			RetryCount *yamlnode.WithDataType[int]   `yaml:"retry_count,omitempty"`
		}]
	}

	var conf Configuration
	err := yaml.Unmarshal([]byte(in), &conf)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test", conf.Data.Database)
	assert.Equal(t, "27017", conf.Data.Backend.Data.Port())
	assert.Equal(t, 5, conf.Data.RetryCount.Data)
	assert.Equal(t, yaml.ScalarNode, conf.Data.Backend.Node.Kind)
	assert.Equal(t, yaml.MappingNode, conf.Node.Kind)
	assert.Equal(t, yaml.ScalarNode, conf.Data.RetryCount.Node.Kind)

	conf.Data.Database = "example"
	conf.Data.RetryCount = nil

	updatedBuf, err := yaml.Marshal(conf)
	if err != nil {
		t.Fatal(err)
	}

	var updated Configuration
	err = yaml.Unmarshal(updatedBuf, &updated)
	require.NoError(t, err)
	assert.Equal(t, "example", updated.Data.Database)
	assert.Nil(t, updated.Data.RetryCount)
}

type Server struct {
	*url.URL
}

func (u *Server) UnmarshalYAML(node *yaml.Node) error {
	var err error
	u.URL, err = url.Parse(node.Value)
	return err
}

func (u *Server) MarshalYAML() (any, error) {
	if u == nil || u.URL == nil {
		return "", nil
	}
	return u.String(), nil
}
