package reader

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func MatchFileExt(ext string) func(os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		return filepath.Ext(info.Name()) == ext
	}
}

func FilterSmallThan(size int64) func(os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		return info.Size() <= size
	}
}

func FilterBigThan(size int64) func(os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		return info.Size() >= size
	}
}

func FilterList(fss ...func(os.FileInfo) bool) func(os.FileInfo) bool {
	return func(info os.FileInfo) bool {
		for _, f := range fss {
			if !f(info) {
				return false
			}
		}
		return true
	}
}

// recurse
func Files(dir string, matchs ...func(os.FileInfo) bool) []string {
	match := FilterList(matchs...)

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

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
