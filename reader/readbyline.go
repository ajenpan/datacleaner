package reader

import (
	"bufio"
	"os"
)

func NewByLine(filename string) (*ByLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
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
	return nil, r.scanner.Err()
}
