package utils

import (
	"log"
	"os"
)

// ReadFileContent reads the content of a file.
func ReadFileContent(filename string) ([]byte, error) {
	dir := os.Args[2]
	filePath := dir + filename

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("failed to read file %s: %v", filePath, err)
		return nil, err
	}

	return fileContent, nil
}

// WriteFileContent writes content to a file.
func WriteFileContent(filename string, content string) error {
	dir := os.Args[2]
	filePath := dir + filename

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Printf("failed to write file %s: %v", filePath, err)
		return err
	}

	return nil
}
