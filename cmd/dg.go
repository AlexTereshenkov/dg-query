/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"os"
)

var ReadFile = func(filePath string) []byte {
	jsonData, readingFileError := os.ReadFile(filePath)
	if readingFileError != nil {
		panic(readingFileError)
	}
	return jsonData
}

func loadJsonFile(jsonData []byte) map[string][]string {
	var adjacencyList map[string][]string
	loadingJsonError := json.Unmarshal(jsonData, &adjacencyList)
	if loadingJsonError != nil {
		panic(loadingJsonError)
	}
	return adjacencyList
}
