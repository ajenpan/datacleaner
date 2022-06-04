package reader

import (
	"bufio"
	"math/rand"
	"os"
	"syscall"
	"time"
)

type FileInfo struct {
	RandomLines []string
	FilePath    string
	FileSize    int64
	CreateTime  time.Time
}

func ReadFileInfo(fileName string) (*FileInfo, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// 1000 line pick out 20 random line
	temp := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		if !scanner.Scan() {
			break
		}
		temp = append(temp, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())

	pickCount := 20
	if len(temp) < pickCount {
		pickCount = len(temp)
	}

	head := make([]string, 0, pickCount)
	for i := 0; i < pickCount; i++ {
		head = append(head, temp[rand.Intn(len(temp))])
	}

	fs, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}

	stat := fs.Sys().(*syscall.Win32FileAttributeData)
	t := time.Unix(0, stat.CreationTime.Nanoseconds())

	ret := &FileInfo{
		RandomLines: head,
		FilePath:    fileName,
		CreateTime:  t,
		FileSize:    fs.Size(),
	}

	return ret, nil
}
