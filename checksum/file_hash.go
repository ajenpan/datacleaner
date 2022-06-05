package checksum

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"os"
)

// MD5File returns MD5 checksum of filename
func MD5File(filename string) (string, error) {
	return sumFile(md5.New(), filename)
}

// SHA256File returns SHA256 checksum of filename
func SHA256File(filename string) (string, error) {
	return sumFile(sha256.New(), filename)
}

// SHA1File returns SHA1 checksum of filename
func SHA1File(filename string) (string, error) {
	return sumFile(sha1.New(), filename)
}

// sumFile calculates the hash based on a provided hash provider
func sumFile(hashAlgorithm hash.Hash, filename string) (string, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return sumReader(hashAlgorithm, bufio.NewReader(file))
}
