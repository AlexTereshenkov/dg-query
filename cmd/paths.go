/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

func paths(filePath string, fromTarget string, toTarget string, readFile ReadFileFunc) ([][]string, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}

	// Initialize the result slice to store all paths
	var result [][]string
	//memo := make(map[string][][]string)
	dfsWithMemoization(adjacencyList, fromTarget, toTarget, fromTarget, []string{}, &result, make(map[string]bool))

	return result, nil
}

// dfs does a depth-first search (DFS) to find all paths from the "from" target to the "to" target
func dfsWithMemoization(
	adjList map[string][]string,
	currentNode,
	endNode,
	startNode string,
	currentPath []string,
	result *[][]string,
	visited map[string]bool,
) {
	// Add the current node to the path
	currentPath = append(currentPath, currentNode)

	// If we reach the end node, add the current path to the result
	if currentNode == endNode {
		*result = append(*result, append([]string{}, currentPath...)) // Append a copy of the currentPath
		// Backtrack by removing the current node from the path
		return
	}

	// Mark the current node as visited
	visited[currentNode] = true

	// Continue exploring neighbors
	for _, neighbor := range adjList[currentNode] {
		// Skip the neighbor if it is already visited or if it's the start node
		if visited[neighbor] || (neighbor == startNode && currentNode != startNode) {
			continue
		}

		// Recursively visit neighbors
		dfsWithMemoization(adjList, neighbor, endNode, startNode, currentPath, result, visited)
	}

	// Backtrack: unmark the current node as visited for other paths
	visited[currentNode] = false
	// Not needed?! Remove the current node from the path for backtracking
	// currentPath = currentPath[:len(currentPath)-1] // Backtrack to remove the current node from the path
}
