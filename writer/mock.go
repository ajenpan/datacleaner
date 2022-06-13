package writer

import (
	"encoding/json"
	"fmt"

	"datacleaner/object"
)

var Printer = func(in object.Object) {
	raw, _ := json.Marshal(in)
	fmt.Println(string(raw))
}

func NewMock(...interface{}) (*Mock, error) {
	return &Mock{
		Printer: Printer,
	}, nil
}

type Mock struct {
	Printer func(line map[string]interface{})
}

func (m *Mock) Write(line object.Object) error {
	if m.Printer != nil {
		m.Printer(line)
	}
	return nil
}

func (m *Mock) Close() error {
	return nil
}
