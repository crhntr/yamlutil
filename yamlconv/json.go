package yamlconv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	nullTag   = "!!null"
	strTag    = "!!str"
	intTag    = "!!int"
	floatTag  = "!!float"
	binaryTag = "!!binary"
)

func ToJSON(buf []byte, node *yaml.Node) ([]byte, error) {
	if node == nil {
		node = &yaml.Node{}
	}
	switch node.Kind {
	case yaml.DocumentNode, 0:
		if len(node.Content) == 0 ||
			(len(node.Content) == 1 && node.Content[0].Tag == nullTag) {
			buf = append(buf, []byte(`{}`)...)
			return buf, nil
		}
		return ToJSON(buf, node.Content[0])
	case yaml.AliasNode:
		return ToJSON(buf, node.Alias)
	case yaml.MappingNode:
		buf = append(buf, '{')
		var err error
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i+1]
			if i < len(node.Content)-1 {
				buf = append(buf, strconv.Quote(key.Value)...)
			}
			buf = append(buf, ':')
			buf, err = ToJSON(buf, value)
			if err != nil {
				return nil, err
			}
			if i < len(node.Content)-2 {
				buf = append(buf, ',')
			}
		}
		buf = append(buf, '}')
		return buf, nil
	case yaml.ScalarNode:
		return encodeScalar(buf, node)
	case yaml.SequenceNode:
		buf = append(buf, '[')
		var err error
		for i, element := range node.Content {
			buf, err = ToJSON(buf, element)
			if err != nil {
				return nil, err
			}
			if i < len(node.Content)-1 {
				buf = append(buf, ',')
			}
		}
		buf = append(buf, ']')
		return buf, nil
	default:
		panic(fmt.Errorf("unknown YAML node kind %d", node.Kind))
	}
}

func encodeScalar(buf []byte, node *yaml.Node) ([]byte, error) {
	switch node.Tag {
	case strTag, binaryTag:
		valueBuf, err := json.Marshal(node.Value)
		return append(buf, valueBuf...), err
	case intTag:
		buf = append(buf, node.Value...)
		return buf, nil
	case floatTag:
		var val float64
		err := node.Decode(&val)
		if err != nil {
			return nil, err
		}
		valBuf, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		addDotZero := false
		if i := strings.IndexByte(node.Value, '.'); i > 0 {
			addDotZero = true
			for _, c := range node.Value[i+1:] {
				if c != '0' {
					addDotZero = false
					break
				}
			}
		}
		if addDotZero {
			valBuf = append(valBuf, ".0"...)
		}
		return append(buf, valBuf...), nil
	case nullTag:
		buf = append(buf, "null"...)
		return buf, nil
	default:
		var v any
		err := node.Decode(&v)
		if err != nil {
			return nil, err
		}
		valueBuf, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buf = append(buf, valueBuf...)
		return buf, nil
	}
}
