/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"fmt"
	"slices"

	"strings"

	"github.com/spf13/cobra"
)

/*
List dependencies for given targets.
*/
func dependencies(cmd *cobra.Command, targets []string) {
	filePath, _ := cmd.Flags().GetString("dg")
	jsonData := ReadFile(filePath)
	adjacencyList := loadJsonFile(jsonData)
	var deps []string

	transitive, _ := cmd.Flags().GetBool("transitive")
	if transitive {
		deps = getDepsTransitive(adjacencyList, targets)
	} else {
		deps = getDepsDirect(adjacencyList, targets)
	}

	reflexive, _ := cmd.Flags().GetBool("reflexive")
	if reflexive {
		reflexiveTargets := getReflexiveTargets(targets, adjacencyList)
		deps = append(deps, reflexiveTargets...)
	}
	slices.Sort(deps)
	deps = slices.Compact(deps)
	output := strings.Join(deps, "\n")
	fmt.Fprintln(cmd.OutOrStdout(), output)

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
	var deps []string
	for _, target := range targets {
		deps = append(deps, adjacencyList[target]...)
	}
	return deps
}

func getDepsTransitive(adjacencyList map[string][]string, targets []string) []string {
	var deps []string
	// Keeping track of visited targets to skip duplicates and handle infinite loops
	visited := make(map[string]bool)

	// Necessary to declare beforehand since it's called recursively
	var getDeps func(target string)
	getDeps = func(target string) {
		// Visited targets can be skipped to avoid infinite recursion
		if visited[target] {
			return
		}
		visited[target] = true

		for _, dep := range adjacencyList[target] {
			deps = append(deps, dep)

			// If the dependency is also a key in the adjacency list, recurse into it
			if _, isKey := adjacencyList[dep]; isKey {
				getDeps(dep)
			}
		}
	}

	// Get transitive dependencies for each target
	for _, target := range targets {
		getDeps(target)
	}
	return deps
}
