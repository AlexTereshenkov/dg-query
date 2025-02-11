/*
Copyright © 2025 Alexey Tereshenkov
*/
package cmd

// to be used in non-unit tests
var ExtractSubgraph = extractSubgraph

// ExtractDependencySubgraph returns a new subgraph as adjacency list containing only the nodes reachable from the given root node
func extractSubgraph(filePath string, rootNode string, readFile ReadFileFunc) (AdjacencyList, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}

	if _, exists := adjacencyList[rootNode]; !exists {
		return AdjacencyList{}, nil
	}

	result := make(AdjacencyList)
	visited := make(map[string]bool)

	var dfs func(node string)
	dfs = func(node string) {
		if visited[node] || adjacencyList[node] == nil {
			return
		}

		visited[node] = true

		result[node] = make([]string, len(adjacencyList[node]))
		copy(result[node], adjacencyList[node])

		for _, neighbor := range adjacencyList[node] {
			dfs(neighbor)
		}
	}

	dfs(rootNode)
	return result, nil
}
