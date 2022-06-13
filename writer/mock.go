package writer

func NewMock(...interface{}) (*Mock, error) {
	return &Mock{}, nil
}

type Mock struct {
	Printer func(line map[string]interface{})
}

func (m *Mock) Write(line map[string]interface{}) error {
	if m.Printer != nil {
		m.Printer(line)
	}
	return nil
}

func (m *Mock) Close() error {
	return nil
}
