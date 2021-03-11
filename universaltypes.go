package sflags

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// -- stringMapValue
type stringMapValue struct {
	value reflect.Value
}

var _ RepeatableFlag = (*stringMapValue)(nil)
var _ Value = (*stringMapValue)(nil)
var _ Getter = (*stringMapValue)(nil)

func newStringMapValue(m reflect.Value) *stringMapValue {
	return &stringMapValue{
		value: m,
	}
}

func (v *stringMapValue) Set(s string) error {
	ss := strings.SplitN(s, ":", 2)
	if len(ss) < 2 {
		return errors.New("invalid map flag syntax, use -map=key1:val1")
	}

	key := ss[0]

	valStr := ss[1]

	val := reflect.New(v.value.Type().Elem()).Interface()
	if err := json.Unmarshal([]byte(valStr), &val); err != nil {
		return err
	}
	v.value.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val).Elem())

	return nil
}

func (v *stringMapValue) Get() interface{} {
	if v != nil &&
		v.value.IsValid() &&
		v.value.Interface() != nil {
		return v.value.Interface()
	}
	return nil
}

func (v *stringMapValue) String() string {
	if v != nil &&
		v.value.IsValid() &&
		v.value.Interface() != nil &&
		v.value.Len() > 0 {
		return fmt.Sprintf("%v", v.value.String())
	}
	return ""
}

func (v *stringMapValue) Type() string { return "map[string]interface{}" }

func (v *stringMapValue) IsCumulative() bool {
	return true
}

// -- sliceValue
type sliceValue struct {
	value reflect.Value
}

var _ RepeatableFlag = (*sliceValue)(nil)
var _ Value = (*sliceValue)(nil)
var _ Getter = (*sliceValue)(nil)

func newSliceValue(slice reflect.Value) *sliceValue {
	return &sliceValue{
		value: slice,
	}
}

func (v *sliceValue) Set(raw string) error {
	val := reflect.New(v.value.Type()).Interface()
	if err := json.Unmarshal([]byte(raw), &val); err != nil {
		return err
	}
	v.value.Set(reflect.ValueOf(val).Elem())
	return nil
}

func (v *sliceValue) Get() interface{} {
	if v != nil &&
		v.value.IsValid() &&
		v.value.Interface() != nil {
		return v.value.Interface()
	}
	return ([]interface{})(nil)
}

func (v *sliceValue) String() string {
	if v == nil ||
		!v.value.IsValid() ||
		v.value.Interface() == nil ||
		v.value.Len() == 0 {
		return "[]"
	}
	out := make([]string, 0, v.value.Len())
	for i := 0; i < v.value.Len(); i++ {
		elem := v.value.Index(i).Interface()
		v, _ := json.Marshal(elem)
		out = append(out, string(v))
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (v *sliceValue) Type() string { return "Slice" }

func (v *sliceValue) IsCumulative() bool {
	return false
}
