package writer

import "datacleaner/object"

type Writer interface {
	Write(object.Object) error
	Close() error
}
