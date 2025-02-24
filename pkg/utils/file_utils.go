package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetAbsolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func ResolveRelativePathToAbsolute(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filename, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}

	absPath := filepath.Join(wd, filename)

	return absPath, nil
}

func GetAbsolutePathDir(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filepath.Dir(filename), nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}

	absPath := filepath.Join(wd, filename)

	return filepath.Dir(absPath), nil
}

func GetFileNameAndType(filePath string) (name, _type string) {
	baseName := filepath.Base(filePath)

	fileType := filepath.Ext(baseName)
	fileType = strings.TrimPrefix(fileType, ".")

	fileName := strings.TrimSuffix(baseName, fileType)
	fileName = strings.TrimSuffix(fileName, ".")

	return fileName, fileType
}

func IsSimpleFileName(fileName string) bool {
	if strings.HasPrefix(fileName, "/") {
		return false
	}

	if strings.Contains(fileName, "/") {
		return false
	}

	return true
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}
