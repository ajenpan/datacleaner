package writer

type Writer interface {
	Write(line map[string]interface{}) error
	Close() error
}
