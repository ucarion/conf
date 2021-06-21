package conf

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagRedact = "redact"
)

// Redact returns a deep copy of v, excluding fields tagged with "redact" set to
// a truthy value.
//
// A truthy value is any value which strconv.Bool parses as "true". If a field
// is tagged with "redact" set to a value unsupported by strconv.Bool, Redact
// will panic.
//
// v must be a struct. This is in contrast to Load, which requires a pointer to
// a struct. If v is not a struct, Redact will panic.
func Redact(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("conf: Redact called on %v (only structs are acceptable)", val.Kind()))
	}

	return redact(val).Interface()
}

func redact(v reflect.Value) reflect.Value {
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return v
	}

	out := reflect.New(t)
	for i := 0; i < t.NumField(); i++ {
		if shouldRedact(t.Field(i)) || !out.Elem().Field(i).CanSet() {
			continue
		}

		out.Elem().Field(i).Set(redact(v.Field(i)))
	}

	return out.Elem()
}

func shouldRedact(f reflect.StructField) bool {
	t, ok := f.Tag.Lookup(tagRedact)
	if !ok {
		return false
	}

	b, err := strconv.ParseBool(t)
	if err != nil {
		panic(fmt.Errorf("conf: error parsing redact tag: %v", err))
	}

	return b
}
