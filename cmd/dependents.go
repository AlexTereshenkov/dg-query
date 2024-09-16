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
Reverse adjacency lists:
{"foo": ["bar", "baz"], "spam": ["eggs", "bar"]} ->
{"bar": ["foo", "spam"], "baz":
}
*/
func reverseAdjacencyLists(adjacencyList AdjacencyList) AdjacencyList {
	reversed := make(AdjacencyList)
	for node, dependencies := range adjacencyList {
		for _, dep := range dependencies {
			reversed[dep] = append(reversed[dep], node)
		}
	}
	return reversed
}

/*
List dependents for given targets.
*/
func dependents(cmd *cobra.Command, targets []string) {
	var adjacencyList AdjacencyList

	dgFilePath, _ := cmd.Flags().GetString("dg")
	rdgFilePath, _ := cmd.Flags().GetString("rdg")
	if rdgFilePath != "" {
		jsonData := ReadFile(rdgFilePath)
		adjacencyList = loadJsonFile(jsonData)
	} else {
		jsonData := ReadFile(dgFilePath)
		adjacencyList = reverseAdjacencyLists(loadJsonFile(jsonData))
	}

	var rdeps []string
	transitive, _ := cmd.Flags().GetBool("transitive")
	if transitive {
		rdeps = getDepsTransitive(adjacencyList, targets)
	} else {
		rdeps = getDepsDirect(adjacencyList, targets)
	}

	reflexive, _ := cmd.Flags().GetBool("reflexive")
	if reflexive {
		reflexiveTargets := getReflexiveTargets(targets, adjacencyList)
		rdeps = append(rdeps, reflexiveTargets...)
	}
	slices.Sort(rdeps)
	rdeps = slices.Compact(rdeps)
	output := strings.Join(rdeps, "\n")
	fmt.Fprintln(cmd.OutOrStdout(), output)

}
