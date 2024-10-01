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

	dfs(adjacencyList, fromTarget, toTarget, []string{fromTarget}, &result)

	return result, nil
}

// dfs does a depth-first search (DFS) to find all paths from the "from" target to the "to" target
func dfs(adjacencyList AdjacencyList, currentNode string, endNode string, currentPath []string, result *[][]string) {
	// If we reach the end node, save the current path to the result
	if currentNode == endNode {
		*result = append(*result, append([]string(nil), currentPath...))
		return
	}

	if nextNodes, exists := adjacencyList[currentNode]; exists {
		for _, nextNode := range nextNodes {
			newPath := append(currentPath, nextNode)
			dfs(adjacencyList, nextNode, endNode, newPath, result)
		}
	}
}
