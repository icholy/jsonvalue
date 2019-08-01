package jsonvalue

import (
	"testing"

	"gotest.tools/assert"
)

var data = []byte(`

	{
		"foo": 123,
		"bar": [1, 2, 3]
	}

`)

func TestValueKey(t *testing.T) {
	v := Parse(data).Key("foo")
	assert.NilError(t, v.Err)
	assert.DeepEqual(t, v.Path, []string{"foo"})
	assert.DeepEqual(t, v.Value, float64(123))
}

func TestValueKeyMissing(t *testing.T) {
	v := Parse(data).Key("non_existent_key")
	assert.Error(t, v.Err, `key not found "non_existent_key"`)
	assert.Check(t, len(v.Path) == 0)
}

func TestValueKeyAndIndex(t *testing.T) {
	v := Parse(data).Key("bar").Index(1)
	assert.NilError(t, v.Err)
	assert.DeepEqual(t, v.Value, float64(2))
	assert.DeepEqual(t, v.Path, []string{"bar", "1"})
}
