package reader

import (
	"bufio"
	"os"
)

type ByLine struct {
	filename string
}

func NewByLine(filename string) *ByLine {
	return &ByLine{filename}
}

func (r *ByLine) Run(c chan string) error {
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
