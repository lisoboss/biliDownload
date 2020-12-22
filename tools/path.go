package tools

import "os"

func CreateDir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	return err
}

func CreateDirFromFilePath(filePath string) error {
	file, err := os.Create(filePath)
	if err == nil {
		_ = file.Close()
		err = os.Remove(filePath)
		return err
	}
	err = CreateDir(filePath)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	return err
}
