/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import "sort"

// to be used in non-unit tests
var Leaves = leaves

// leaves returns nodes (aka sinks) that have no dependencies
func leaves(filePath string, readFile ReadFileFunc) ([]string, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}

	// put all nodes into a single map first to iterate all at once
	nodes := make(map[string]bool)

	for node := range adjacencyList {
		nodes[node] = true
	}
	for _, dependencies := range adjacencyList {
		for _, dep := range dependencies {
			nodes[dep] = true
		}
	}

	result := []string{}
	for node := range nodes {
		dependencies, exists := adjacencyList[node]
		if !exists || len(dependencies) == 0 {
			result = append(result, node)
		}
	}
	sort.Strings(result)
	return result, nil
}
