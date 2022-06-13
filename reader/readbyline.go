package reader

import (
	"bufio"
	"io"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func NewByLine(filename string) (*ByLine, error) {
	return NewByLineWithDecoder(filename, nil)
}

func NewByLineWithDecoder(filename string, decoder *encoding.Decoder) (*ByLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	var r io.Reader = file
	if decoder != nil {
		r = transform.NewReader(r, decoder)
	}
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	return &ByLine{file: file, scanner: scanner}, nil
}

type ByLine struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (r *ByLine) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func (r *ByLine) Read() ([]byte, error) {
	if r.scanner.Scan() {
		return r.scanner.Bytes(), nil
	}
	if r.scanner.Err() != nil {
		return nil, r.scanner.Err()
	}
	return nil, io.EOF
}
