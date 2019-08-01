package jsonvalue

import (
	"testing"

	"gotest.tools/assert"
)

var data = []byte(`

	{
		"foo": 123,
		"bar": [1, 2, 3],
		"baz": null,
		"poo": false
	}

`)

func TestKey(t *testing.T) {
	v := Parse(data).Key("foo")
	assert.NilError(t, v.Err)
	assert.DeepEqual(t, v.Path, []string{"foo"})
	assert.DeepEqual(t, v.Value, float64(123))
}

func TestKeyMissing(t *testing.T) {
	v := Parse(data).Key("non_existent_key")
	assert.Error(t, v.Err, `key not found "non_existent_key"`)
	assert.Check(t, len(v.Path) == 0)
}

func TestKeyAndIndex(t *testing.T) {
	v := Parse(data).Key("bar").Index(1)
	assert.NilError(t, v.Err)
	assert.DeepEqual(t, v.Value, float64(2))
	assert.DeepEqual(t, v.Path, []string{"bar", "1"})
}

func TestIsNull(t *testing.T) {
	v := Parse(data).Key("baz")
	isNull, err := v.IsNull()
	assert.NilError(t, err)
	assert.Assert(t, isNull)
}

func TestType(t *testing.T) {
	v := Parse(data)
	assert.Equal(t, TypeObject, v.Type())
	assert.Equal(t, TypeArray, v.Key("bar").Type())
	assert.Equal(t, TypeNum, v.Key("bar").Index(0).Type())
	assert.Equal(t, TypeNull, v.Key("baz").Type())
	assert.Equal(t, TypeBool, v.Key("poo").Type())
}
