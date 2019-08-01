# JSON Value 

> This package provides a Value type for getting values out of arbitrary json structure


``` go
val := jsonvalue.Parse([]byte(`{ "foo": { "bar": [123, 324, 1] } }`))
num, _ := val.Lookup("foo", "bar").Index(1).Num()
fmt.Println(num)
```
