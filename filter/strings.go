package filter

import (
	"strings"

	"datacleaner/object"
)

type SplitBy struct {
	By   string
	Slot []string

	Field string
	// Name string
}

// func (s *SpliterBy) String() string {
// 	return s.Name
// }

func (s *SplitBy) Do(in object.Object) (object.Object, bool) {
	field := s.Field
	if field == "" {
		field = "_raw"
	}

	res, ok := s.Parse(in[field])

	if !ok {
		return in, false
	}

	for k, v := range res {
		in[k] = v
	}

	return in, true
}

func (s *SplitBy) Parse(raw interface{}) (map[string]interface{}, bool) {
	if raw == nil {
		return nil, false
	}

	line, ok := raw.(string)
	if !ok {
		return nil, false
	}

	if len(line) < len(s.By)*(len(s.Slot)-1) {
		return nil, false
	}

	temp := strings.Split(line, s.By)
	if len(temp) < len(s.Slot) {
		return nil, false
	}

	res := make([]string, 0, len(s.Slot))
	for _, v := range temp {
		vv := strings.TrimSpace(v)
		res = append(res, vv)
		// if len(vv) > 0 {
		// }
	}

	if len(res) != len(s.Slot) {
		return nil, false
	}

	resMap := make(map[string]interface{}, len(s.Slot))
	for i, v := range s.Slot {
		resMap[v] = string(res[i])
	}
	return resMap, true
}

func (s *SplitBy) String() string {
	return "SplitBy"
}
