/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"os"
)

type AdjacencyList map[string][]string

var ReadFile = func(filePath string) []byte {
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
