package filter

type Filter interface {
	Do(in *Element) (*Element, bool)
}

type Element struct {
	Data map[string]interface{}
}

func NewElement(raw string) *Element {
	data := make(map[string]interface{})
	data["_raw"] = raw
	return &Element{Data: data}
}

func NewMultiple(filters ...Filter) Filter {
	return &Multiple{Filters: filters}
}

type Multiple struct {
	Filters []Filter
}

func (f *Multiple) Do(in *Element) (*Element, bool) {
	for _, filter := range f.Filters {
		if in, ok := filter.Do(in); !ok {
			return in, false
		}
	}
	return in, true
}

type Drop struct {
	Field     string
	Condition func(interface{}) bool
}

func (f *Drop) Do(in *Element) (*Element, bool) {
	if f.Condition(in.Data[f.Field]) {
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
	C func(in *Element) (*Element, bool)
}

func (f *Custom) Do(in *Element) (*Element, bool) {
	return f.C(in)
}

func NewCustom(f func(in *Element) (*Element, bool)) Filter {
	return &Custom{C: f}
}

type FieldsDelete struct {
	Fields []string
}

func (f *FieldsDelete) Do(in *Element) (*Element, bool) {
	for _, v := range f.Fields {
		delete(in.Data, v)
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

func (f *FieldsRetain) Do(in *Element) (*Element, bool) {
	for k := range in.Data {
		if _, ok := f.Fields[k]; !ok {
			delete(in.Data, k)
		}
	}
	return in, true
}
