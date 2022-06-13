package main

import (
	"fmt"

	"datacleaner/reader"
)

func main() {

	reader.NewSqlserverReader("")
	fmt.Println("Hello, World!")
	// Shenfengzheng()
	return

	// dir := "D:/passwd_bad/"
	// bigfiles := reader.Files(dir, reader.FilterBigThan(2*1024*1024*1024))
	// for _, v := range bigfiles {
	// 	fmt.Println("big file pass: ", v)
	// }
	// files := reader.Files(dir)
	// utility.RemoveSameFile(files)
	// return

	target := "D:/passwd_data/"

	targets := reader.Files(target, reader.MatchFileExt(".txt"))

	job, err := NewJob("passwd_data_all", targets)
	if err != nil {
		fmt.Println(err)
	}
	err = job.Run()

	if err != nil {
		fmt.Println(err)
	}
}
