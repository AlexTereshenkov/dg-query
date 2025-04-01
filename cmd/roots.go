/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import "sort"

// to be used in non-unit tests
var Roots = roots

// roots returns nodes (aka sources) that no other node depends on
func roots(filePath string, readFile ReadFileFunc) ([]string, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}

	hasDependents := make(map[string]bool)

	for _, dependents := range adjacencyList {
		for _, dependent := range dependents {
			hasDependents[dependent] = true
		}
	}

	result := []string{}
	for node := range adjacencyList {
		if !hasDependents[node] {
			result = append(result, node)
		}
	}

	sort.Strings(result)
	return result, nil
}
