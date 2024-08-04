package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readFileText(filePath string) (string, error) {
	var fileText string

	file, err := os.Open(filePath)
	if err != nil {
		return fileText, fmt.Errorf("Failed to open file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	fileText = strings.Join(fileLines, "\n")

	return fileText, nil
}

