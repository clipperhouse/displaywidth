package testdata

import (
	"crypto/rand"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"runtime"
)

func InvalidUTF8() ([]byte, error) {
	return load("UTF-8-test.txt")
}

func Sample() ([]byte, error) {
	return load("sample.txt")
}

func TestCases() ([]byte, error) {
	return load("test_cases.txt")
}

func RandomBytes() ([]byte, error) {
	length := mathrand.Intn(50)
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func load(filename string) ([]byte, error) {
	// Get the directory of this source file
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	path := filepath.Join(dir, filename)

	return os.ReadFile(path)
}
