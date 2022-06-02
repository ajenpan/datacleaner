package reader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MatchFileExt(ext string) func(os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		return info.Name()[len(info.Name())-len(ext):] == ext
	}
}

// recurse
func WalkDirFiles(dir string, match func(os.FileInfo) bool) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println(path)
		if match != nil {
			if match(info) {
				files = append(files, path)
			}
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// not recurse
func ReadDir(dir string, match func(os.FileInfo) bool) []string {
	finfos, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	files := []string{}

	for _, info := range finfos {
		if info.IsDir() {
			continue
		}

		if match != nil {
			if match(info) {
				files = append(files, filepath.Join(dir, info.Name()))
			}
		} else {
			files = append(files, filepath.Join(dir, info.Name()))
		}
	}

	return files
}
