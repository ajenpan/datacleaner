package reader

import (
	"io"

	"datacleaner/object"
)

func NewMulti(list ...Reader) Reader {
	return &Multi{list: list}
}

type Multi struct {
	list []Reader
}

func (r *Multi) Read() (object.Object, error) {
	if len(r.list) == 0 {
		return nil, io.EOF
	}
	o, err := r.list[0].Read()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}

		r.list[0].Close()
		r.list = r.list[1:]
		return r.Read()
	}
	return o, nil
}

func (r *Multi) Close() error {
	return nil
}
