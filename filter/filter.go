package filter

import (
	"datacleaner/object"
)

type Filter interface {
	Do(in object.Object) (object.Object, bool)
}

func NewMultiple(filters ...Filter) Filter {
	return &Multiple{Filters: filters}
}

type Multiple struct {
	Filters []Filter
}

func (f *Multiple) Do(in object.Object) (object.Object, bool) {
	var ok bool
	for _, filter := range f.Filters {
		if in, ok = filter.Do(in); !ok {
			return in, false
		}
	}
	return in, true
}

type Drop struct {
	Field     string
	Condition func(interface{}) bool
}

func (f *Drop) Do(in object.Object) (object.Object, bool) {
	if f.Condition(in[f.Field]) {
		return in, false
	}
	return in, true
}

func Equal(target interface{}) func(interface{}) bool {
	return func(v interface{}) bool {
		return v == target
	}
}

func NotEqual(target interface{}) func(interface{}) bool {
	return func(v interface{}) bool {
		return v != target
	}
}

type Custom struct {
	C func(in object.Object) (object.Object, bool)
}

func (f *Custom) Do(in object.Object) (object.Object, bool) {
	return f.C(in)
}

func NewCustom(f func(in object.Object) (object.Object, bool)) Filter {
	return &Custom{C: f}
}

type FieldsDelete struct {
	Fields []string
}

func (f *FieldsDelete) Do(in object.Object) (object.Object, bool) {
	for _, v := range f.Fields {
		delete(in, v)
	}
	return in, true
}

func NewFieldsRetain(fields ...string) *FieldsRetain {
	m := make(map[string]struct{}, len(fields))
	for _, v := range fields {
		m[v] = struct{}{}
	}
	return &FieldsRetain{Fields: m}
}

type FieldsRetain struct {
	Fields map[string]struct{}
}

func (f *FieldsRetain) Do(in object.Object) (object.Object, bool) {
	for k := range in {
		if _, ok := f.Fields[k]; !ok {
			delete(in, k)
		}
	}
	return in, true
}
