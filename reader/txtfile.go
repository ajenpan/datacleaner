package reader

import (
	"bufio"
	"os"
)

type TxtFile struct {
	filename string
}

func NewTxtFile(filename string) *TxtFile {
	return &TxtFile{filename}
}

func (r *TxtFile) Run(c chan string) error {
	file, err := os.Open(r.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		c <- scanner.Text()
	}
	return nil
}
