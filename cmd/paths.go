/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

// to be used in non-unit tests
var Paths = paths

func paths(filePath string, fromTarget string, toTarget string, maxPaths int, readFile ReadFileFunc) ([][]string, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}
	var result [][]string
	dfsWithMemoization(adjacencyList, fromTarget, toTarget, fromTarget, []string{}, &result, make(map[string]bool), maxPaths)

	return result, nil
}

// dfs does a depth-first search (DFS) to find all (or some) paths from the "from" target to the "to" target
func dfsWithMemoization(
	adjacencyList map[string][]string,
	currentNode,
	endNode,
	startNode string,
	currentPath []string,
	result *[][]string,
	visited map[string]bool,
	maxPaths int,
) {
	if maxPaths > 0 && len(*result) >= maxPaths {
		return
	}

	// Add the current node to the path
	currentPath = append(currentPath, currentNode)

	// If we reach the end node, add the current path to the result
	if currentNode == endNode {
		*result = append(*result, append([]string{}, currentPath...))
		// Backtrack by removing the current node from the path
		return
	}

	visited[currentNode] = true

	// Continue exploring neighbors
	for _, neighbor := range adjacencyList[currentNode] {
		// Skip the neighbor if it has already been visited or if it is the start node
		if visited[neighbor] || (neighbor == startNode && currentNode != startNode) {
			continue
		}

		// Recursively visit neighbors
		dfsWithMemoization(adjacencyList, neighbor, endNode, startNode, currentPath, result, visited, maxPaths)
	}

	// Backtrack: mark the current node as not visited for other paths
	visited[currentNode] = false
}
