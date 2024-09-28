/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"os"
)

// Function type to be used for reading files
type ReadFileFunc func(filePath string) ([]byte, error)

type AdjacencyList map[string][]string

var DefaultReadFile = func(filePath string) ([]byte, error) {
	jsonData, readingFileError := os.ReadFile(filePath)
	if readingFileError != nil {
		return nil, readingFileError
	}
	return jsonData, nil
}

func loadJsonFile(jsonData []byte) (AdjacencyList, error) {
	var adjacencyList AdjacencyList
	loadingJsonError := json.Unmarshal(jsonData, &adjacencyList)
	if loadingJsonError != nil {
		return nil, loadingJsonError
	}
	return adjacencyList, nil
}
