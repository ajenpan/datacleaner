package reader

import "datacleaner/object"

type Reader interface {
	Read() (object.Object, error)
	Close() error
}
