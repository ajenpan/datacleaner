package hashsum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

// MD5File returns MD5 checksum of filename
func MD5File(filename string) (string, error) {
	return HashFileSum(md5.New(), filename)
}

// SHA256File returns SHA256 checksum of filename
func SHA256File(filename string) (string, error) {
	return HashFileSum(sha256.New(), filename)
}

// SHA1File returns SHA1 checksum of filename
func SHA1File(filename string) (string, error) {
	return HashFileSum(sha1.New(), filename)
}

// MD5sumReader returns MD5 checksum of content in reader
func MD5sumReader(reader io.Reader) (string, error) {
	return HashSum(md5.New(), reader)
}

// SHA256sumReader returns SHA256 checksum of content in reader
func SHA256sumReader(reader io.Reader) (string, error) {
	return HashSum(sha256.New(), reader)
}

// SHA1sumReader returns SHA1 checksum of content in reader
func SHA1sumReader(reader io.Reader) (string, error) {
	return HashSum(sha1.New(), reader)
}

// HashFileSum calculates the hash based on a provided hash provider
func HashFileSum(hashAlgorithm hash.Hash, filename string) (string, error) {
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

	return HashSum(hashAlgorithm, file)
}

// HashSum calculates the hash based on a provided hash provider
func HashSum(hashAlgorithm hash.Hash, reader io.Reader) (string, error) {
	var returnSHA1String string

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hashAlgorithm, reader); err != nil {
		return returnSHA1String, err
	}

	hashInBytes := hashAlgorithm.Sum(nil)

	//Convert the bytes to a string
	returnSHA1String = hex.EncodeToString(hashInBytes)
	return returnSHA1String, nil
}
