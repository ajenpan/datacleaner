package reader

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"

	"datacleaner/object"
)

type XLSXSlot struct {
	Key   string
	Index int
}

type XLSXReader struct {
	xlsxFile *excelize.File
	rows     *excelize.Rows
	solts    []*XLSXSlot
}

func NewXLSXReader(file string, solts []*XLSXSlot) (*XLSXReader, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}
	rows, err := f.Rows("Sheet1")
	if err != nil {
		return nil, err
	}
	return &XLSXReader{
		xlsxFile: f,
		rows:     rows,
		solts:    solts,
	}, nil
}

func (r *XLSXReader) Read() (object.Object, error) {
	if !r.rows.Next() {
		return nil, io.EOF
	}
	row, err := r.rows.Columns()
	if err != nil {
		return nil, err
	}

	obj := object.New()
	for _, s := range r.solts {
		if s.Index >= len(row) {
			fmt.Println("index is bigger than row length", s.Index, len(row))
			continue
		}
		obj[s.Key] = row[s.Index]
	}
	return obj, nil
}

func (r *XLSXReader) Close() error {
	r.xlsxFile.Close()
	return nil
}
