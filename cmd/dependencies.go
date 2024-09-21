/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"slices"
)

// to be used in non-unit tests
var Dependencies = dependencies

/*
List dependencies for given targets.
*/
func dependencies(filePath string, targets []string, transitive bool, reflexive bool,
	depth int, readFile ReadFileFunc) []string {
	jsonData := readFile(filePath)
	adjacencyList := loadJsonFile(jsonData)
	var deps []string

	if transitive {
		deps = getDepsTransitive(adjacencyList, targets, depth)
	} else {
		deps = getDepsDirect(adjacencyList, targets)
	}

	if reflexive {
		reflexiveTargets := getReflexiveTargets(targets, adjacencyList)
		deps = append(deps, reflexiveTargets...)
	}
	slices.Sort(deps)
	deps = slices.Compact(deps)
	return deps

}

func getReflexiveTargets(targets []string, adjacencyList map[string][]string) []string {
	var candidates []string
	for _, target := range targets {
		_, exists := adjacencyList[target]
		if exists {
			candidates = append(candidates, target)
		}
	}
	return candidates
}

func getDepsDirect(adjacencyList map[string][]string, targets []string) []string {
	deps := []string{}
	for _, target := range targets {
		deps = append(deps, adjacencyList[target]...)
	}
	return deps
}

func getDepsTransitive(adjacencyList map[string][]string, targets []string, depth int) []string {
	deps := []string{}
	visited := make(map[string]bool)

	// Necessary to declare beforehand since it's called recursively
	var getDeps func(target string, currentDepth int)
	getDeps = func(target string, currentDepth int) {
		// Visited targets can be skipped to avoid infinite recursion
		if visited[target] {
			return
		}
		visited[target] = true

		for _, dep := range adjacencyList[target] {
			deps = append(deps, dep)

			// If the dependency is also a key in the adjacency list, recurse into it
			if _, isKey := adjacencyList[dep]; isKey {
				if depth != 0 && currentDepth >= depth {
					continue
				}
				getDeps(dep, currentDepth+1)
			}
		}
	}

	// Get transitive dependencies for each target
	for _, target := range targets {
		getDeps(target, 1)
	}
	slices.Sort(deps)
	deps = slices.Compact(deps)
	return deps
}
