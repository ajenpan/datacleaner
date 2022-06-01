package main

import (
	"fmt"
	"os"
	"strings"

	"datacleaner/reader"
)

func main() {
	fmt.Println("Hello, World!")

	files := reader.ReadDir("./", func(finfo os.FileInfo) bool {
		return strings.HasSuffix(finfo.Name(), ".go")
	})

	for _, f := range files {
		c := make(chan string, 100)
		go reader.NewTxtFile(f).Run(c)

		for line := range c {
			fmt.Println(line)
		}
	}
}
