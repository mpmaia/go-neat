package test

import (
	"io/ioutil"
	"path/filepath"
)

// GetTempPath Get Temp path
func GetTempPath(fileName string) string {
	dir, err := ioutil.TempDir("", "neatdb-test")
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, fileName)
}
