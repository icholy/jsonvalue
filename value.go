package jsonvalue

import (
	"encoding/json"
	"fmt"
)

// Parse json data into a Value
func Parse(data []byte) Value {
	var v Value
	if err := json.Unmarshal(data, &v.Value); err != nil {
		v.Err = fmt.Errorf("parse: %v", err)
	}
	return v
}

// Value is an object, array, string, number, or bool
type Value struct {
	Err   error
	Path  string
	Value interface{}
}

// String implements fmt.Stringer
func (v Value) String() string {
	if v.Err != nil {
		return fmt.Sprintf("%s: %v", v.Path, v.Err)
	}
	return fmt.Sprint(v.Value)
}

// Walk calls fn with all elements of value recursively.
// When fn returns false, it doesn't recurse further.
func (v Value) Walk(fn func(v Value) bool) {
	if !fn(v) {
		return
	}
	switch v.Value.(type) {
	case map[string]interface{}:
		obj, _ := v.Object()
		for _, v := range obj {
			v.Walk(fn)
		}
	case []interface{}:
		arr, _ := v.Array()
		for _, v := range arr {
			v.Walk(fn)
		}
	}
}

// Str returns the value as a string or an error
// if the value is not a string
func (v Value) Str() (string, error) {
	if v.Err != nil {
		return "", v.Err
	}
	s, ok := v.Value.(string)
	if !ok {
		return "", fmt.Errorf("%s: not a string", v.Path)
	}
	return s, nil
}

// Bool returns the value as a bool or an error
// if the value is not a bool
func (v Value) Bool() (bool, error) {
	if v.Err != nil {
		return false, v.Err
	}
	b, ok := v.Value.(bool)
	if !ok {
		return false, fmt.Errorf("%s: not a boolean", v.Path)
	}
	return b, nil
}

// Num returns the value as  a float64 or a n error
// if the value is not a float64
func (v Value) Num() (float64, error) {
	if v.Err != nil {
		return 0, v.Err
	}
	x, ok := v.Value.(float64)
	if !ok {
		return 0, fmt.Errorf("%s: not a boolean", v.Path)
	}
	return x, nil
}

// Object returns the value as a map of values or an error
// if the value is not an object
func (v Value) Object() (map[string]Value, error) {
	if v.Err != nil {
		return nil, v.Err
	}
	m, ok := v.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s: not an object", v.Path)
	}
	obj := map[string]Value{}
	for k, v0 := range m {
		obj[k] = Value{
			Value: v0,
			Path:  fmt.Sprintf("%s[%q]", v.Path, k),
		}
	}
	return obj, nil
}

// Array returns the value as a slice of values or an error
// if the value is not an array
func (v Value) Array() ([]Value, error) {
	if v.Err != nil {
		return nil, v.Err
	}
	s, ok := v.Value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s: not an array", v.Path)
	}
	arr := make([]Value, len(s))
	for i, v0 := range s {
		arr[i] = Value{
			Value: v0,
			Path:  fmt.Sprintf("%s[%d]", v.Path, i),
		}
	}
	return arr, nil
}

// Len returns the array lenth or an error
// if the value is not an array
func (v Value) Len() (int, error) {
	if v.Err != nil {
		return 0, v.Err
	}
	s, ok := v.Value.([]interface{})
	if !ok {
		return 0, fmt.Errorf("%s: not an array", v.Path)
	}
	return len(s), nil
}

// Key returns the value at the specified key.
// If the value is not an object, the returned value will contain an error
func (v Value) Key(name string) Value {
	if v.Err != nil {
		return v
	}
	path := fmt.Sprintf("%s[%q]", v.Path, name)
	m, ok := v.Value.(map[string]interface{})
	if !ok {
		return Value{
			Err:  fmt.Errorf("not an object"),
			Path: path,
		}
	}
	x, ok := m[name]
	if !ok {
		return Value{
			Err:  fmt.Errorf("key not found %q", name),
			Path: path,
		}
	}
	return Value{
		Value: x,
		Path:  path,
	}
}

// Index returns the value at the specified index.
// If the value is not an array, the returned value will contain an error
func (v Value) Index(i int) Value {
	if v.Err != nil {
		return v
	}
	path := fmt.Sprintf("%s[%d]", v.Path, i)
	s, ok := v.Value.([]interface{})
	if !ok {
		return Value{
			Err:  fmt.Errorf("not an array"),
			Path: path,
		}
	}
	if i < 0 || i >= len(s) {
		return Value{
			Err:  fmt.Errorf("index out of range %d", i),
			Path: path,
		}
	}
	return Value{
		Value: s[i],
		Path:  path,
	}
}
