package tools

import (
	"os"
	"path/filepath"
)

func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

func CreateDirFromFilePath(filePath string) error {
	filePath = filepath.Dir(filePath)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return CreateDir(filePath)
		}
	}
	return err
}
