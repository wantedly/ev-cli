package util

import (
	"os"
	"path"
)

func GetWorkingDirectoryName() string {
	if dir, err := os.Getwd(); err == nil {
		return path.Base(dir)
	}
	return ""
}
