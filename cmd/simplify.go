/*
Copyright © 2025 Alexey Tereshenkov
*/
package cmd

import (
	"log"
	"strings"
)

// to be used in non-unit tests
var SimplifyAdjacencyList = simplifyAdjacencyList

const (
	TechniqueTransitiveReduction = "transitive-reduction"
)

var allowedTechniques = []string{
	TechniqueTransitiveReduction,
}

func isValidTechnique(technique string) bool {
	for _, allowedTechnique := range allowedTechniques {
		if technique == allowedTechnique {
			return true
		}
	}
	return false
}

// simplifyAdjacencyList simplifies the adjacency list by applying given technique
func simplifyAdjacencyList(filePath string, readFile ReadFileFunc, technique string) (AdjacencyList, error) {
	if !isValidTechnique(technique) {
		log.Printf("invalid technique: %s. Allowed techniques are: %s\n", technique, strings.Join(allowedTechniques, ","))
		return make(map[string][]string), nil
	}
	jsonData, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	adjacencyList, err := loadJsonFile(jsonData)
	if err != nil {
		return nil, err
	}
	return transitiveReduction(adjacencyList), nil
}

func transitiveReduction(adjacencyList AdjacencyList) AdjacencyList {
	for node, directDeps := range adjacencyList {
		if len(directDeps) <= 1 {
			continue
		}

		// tracking which dependencies to keep
		newDeps := []string{}

		for _, dep := range directDeps {
			isTransitive := false
			// check if this dependency is a transitive dependency of any other direct dependency
			for _, otherDep := range directDeps {
				if otherDep == dep {
					continue
				}
				// depth of 0 returns all transitive dependencies
				otherTransitive := getDepsTransitive(adjacencyList, []string{otherDep}, 0)

				// if the `dep` is already in the `otherTransitive`, it is redundant
				for _, t := range otherTransitive {
					if t == dep {
						isTransitive = true
						break
					}
				}

				// once we know that `dep` can be reached through `otherDep` (making it a transitive dependency),
				// we don't need to check any other paths
				if isTransitive {
					break
				}
			}

			// keep the dependency only if it's not a transitive dependency through any other direct dependency
			if !isTransitive {
				newDeps = append(newDeps, dep)
			}
		}

		// update the adjacency list in place which is safe to do as the Go spec guarantees that a range loop
		// on a map will iterate over a snapshot of the map's contents at the start of the loop
		adjacencyList[node] = newDeps
	}

	return adjacencyList
}
