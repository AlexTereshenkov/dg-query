/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import "sort"

func cycles(filePath string, readFile ReadFileFunc) ([][]string, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}

	recursionStack := make(map[string]bool)
	cycles := make(map[string][]string)

	// Sort a cycle starting from a specific node
	sortCycle := func(cycle []string) []string {
		// Find the minimum element
		minIdx := 0
		for i := 1; i < len(cycle); i++ {
			if cycle[i] < cycle[minIdx] {
				minIdx = i
			}
		}

		// Rotate the slice to place the minimum element at the start
		result := make([]string, len(cycle))
		for i := 0; i < len(cycle); i++ {
			result[i] = cycle[(i+minIdx)%len(cycle)]
		}
		return result
	}

	// Create a unique key for given cycle
	cycleKey := func(cycle []string) string {
		sortedCycle := sortCycle(cycle)
		result := ""
		for _, node := range sortedCycle {
			result += node
		}
		return result
	}

	var dfs func(node string, path []string)
	dfs = func(node string, path []string) {
		recursionStack[node] = true
		path = append(path, node)

		for _, neighbor := range adjacencyList[node] {
			// Found a cycle
			if recursionStack[neighbor] {
				// Find the start of the cycle in the path
				cycleStart := -1
				for i, n := range path {
					if n == neighbor {
						cycleStart = i
						break
					}
				}

				if cycleStart != -1 {
					cycle := path[cycleStart:]
					sortedCycle := sortCycle(cycle)
					// Use the sorted cycle string as key to avoid duplicates
					key := cycleKey(cycle)
					cycles[key] = sortedCycle
				}
			}

			// Continue exploring if the neighbor isn't in the current path
			// because there may be overlapping cycles (A->B->C->A & B->C->D->A)
			if !recursionStack[neighbor] {
				dfs(neighbor, path)
			}
		}

		recursionStack[node] = false
	}

	for node := range adjacencyList {
		if !recursionStack[node] {
			dfs(node, []string{})
		}
	}

	// Convert cycles map to slice and sort for consistent output
	result := make([][]string, 0, len(cycles))
	for _, cycle := range cycles {
		result = append(result, cycle)
	}

	// Sort the cycles based on their first element for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0]
	})

	return result, nil
}
