package tools

import "os"

func CreateDir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	return err
}

func CreateDirFromFilePath(filePath string) error {
	err := CreateDir(filePath)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	return err
}
