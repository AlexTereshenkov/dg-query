/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"os"
)

// Function type to be used for reading files
type ReadFileFunc func(filePath string) []byte

type AdjacencyList map[string][]string

var DefaultReadFile = func(filePath string) []byte {
	jsonData, readingFileError := os.ReadFile(filePath)
	if readingFileError != nil {
		panic(readingFileError)
	}
	return jsonData
}

func loadJsonFile(jsonData []byte) AdjacencyList {
	var adjacencyList AdjacencyList
	loadingJsonError := json.Unmarshal(jsonData, &adjacencyList)
	if loadingJsonError != nil {
		panic(loadingJsonError)
	}
	return adjacencyList
}
