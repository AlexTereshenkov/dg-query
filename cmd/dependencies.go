/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"os"
	"strings"

	"github.com/spf13/cobra"
)

/*
List dependencies for given targets.
*/
func dependencies(cmd *cobra.Command, targets []string) {
	filePath, _ := cmd.Flags().GetString("dg")
	jsonData := ReadFile(filePath)
	adjacencyList := loadJsonFile(jsonData)

	var deps []string
	for _, target := range targets {
		deps = append(deps, adjacencyList[target]...)
	}
	output := strings.Join(deps, "\n")
	fmt.Fprintln(cmd.OutOrStdout(), output)

}

// TODO: use interfaces instead to make mocking possible under test
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