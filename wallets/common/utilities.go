package common

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// check string is empty
func IsEmptyString(s string) bool {
	return s == ""
}

func WriteDataToFile(filename string, data []byte) error {
	if IsEmptyString(filename) {
		return errors.New("filename invalid input")
	}
	path := filepath.Dir(filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0600)
	}

	return ioutil.WriteFile(filename, data, 0600)
}
