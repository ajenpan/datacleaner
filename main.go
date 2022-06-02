package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	// writer.Test()
	// files := reader.ReadDir("./", func(finfo os.FileInfo) bool {
	// 	return strings.HasSuffix(finfo.Name(), ".go")
	// })

	// for _, f := range files {
	// 	c := make(chan string, 100)
	// 	go reader.NewByLine(f).Run(c)

	// 	for line := range c {
	// 		fmt.Println(line)
	// 	}
	// }
}
