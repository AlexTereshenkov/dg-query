/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import "sort"

// to be used in non-unit tests
var ListConnectedComponents = listConnectedComponents

// listConnectedComponents lists connected components in a graph given a filepath
func listConnectedComponents(filePath string, readFile ReadFileFunc) ([][]string, error) {
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}
	return getConnectedComponents(adjacencyList), nil
}

// getConnectedComponents finds connected components in a graph
func getConnectedComponents(adjacencyList AdjacencyList) [][]string {

	// Convert directed graph to undirected by adding reverse edges;
	// this is necessary so that when node A is connected to B, we could
	// automatically say that B is connected to A.
	undirectedGraph := make(AdjacencyList)

	// Initialize all nodes from the original adjacency list
	for node := range adjacencyList {
		if _, exists := undirectedGraph[node]; !exists {
			undirectedGraph[node] = make([]string, 0)
		}
	}

	// Add connectivity edges in both directions
	for node, neighbors := range adjacencyList {
		for _, neighbor := range neighbors {
			// Add forward edge
			undirectedGraph[node] = append(undirectedGraph[node], neighbor)
			// Initialize neighbor if it does not exist
			if _, exists := undirectedGraph[neighbor]; !exists {
				undirectedGraph[neighbor] = make([]string, 0)
			}
			// Add reverse edge
			undirectedGraph[neighbor] = append(undirectedGraph[neighbor], node)
		}
	}

	visitedNodes := make(map[string]bool)
	// initialize to an empty slice (returned for empty adjacency list)
	connectedComponents := make([][]string, 0)

	var dfs func(node string, component *[]string)
	dfs = func(node string, component *[]string) {
		visitedNodes[node] = true
		*component = append(*component, node)

		for _, neighbor := range undirectedGraph[node] {
			if !visitedNodes[neighbor] {
				dfs(neighbor, component)
			}
		}
	}

	// Sort all nodes for consistent traversal order; this is necessary
	// to ensure the components are reported back in the same order
	nodes := make([]string, 0)
	for node := range undirectedGraph {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	// Find components starting from each unvisited node in the
	// undirected graph that was built earlier
	for _, node := range nodes {
		if !visitedNodes[node] {
			component := make([]string, 0)
			dfs(node, &component)
			sort.Strings(component)
			connectedComponents = append(connectedComponents, component)
		}
	}

	// Sort components by their first element for readability
	sort.Slice(connectedComponents, func(i, j int) bool {
		return connectedComponents[i][0] < connectedComponents[j][0]
	})

	return connectedComponents
}
