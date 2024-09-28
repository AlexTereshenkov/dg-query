/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"slices"
)

/*
Reverse adjacency lists:
{"foo": ["bar", "baz"]} ->
{"bar": ["foo"], "baz": ["foo"]
}
*/
func reverseAdjacencyLists(adjacencyList AdjacencyList) AdjacencyList {
	reversed := make(AdjacencyList)
	for node, dependencies := range adjacencyList {
		for _, dep := range dependencies {
			reversed[dep] = append(reversed[dep], node)
			slices.Sort(reversed[dep])
		}
	}
	return reversed
}

// to be used in non-unit tests
var Dependents = dependents

/*
List dependents for given targets. If the reverse dependency graph is provided,
it's used, otherwise a dependency graph is reversed first.
*/
func dependents(filePathDg string, filePathDgReverse string, targets []string, transitive bool, reflexive bool,
	depth int, DefaultReadFile ReadFileFunc) ([]string, error) {
	var adjacencyList AdjacencyList

	if filePathDgReverse != "" {
		jsonData, err := DefaultReadFile(filePathDgReverse)
		if err != nil {
			return nil, err
		}
		adjacencyList, err = loadJsonFile(jsonData)
		if err != nil {
			return nil, err
		}
	} else {
		jsonData, err := DefaultReadFile(filePathDg)
		if err != nil {
			return nil, err
		}
		adjacencyList, err = loadJsonFile(jsonData)
		if err != nil {
			return nil, err
		}
		adjacencyList = reverseAdjacencyLists(adjacencyList)
	}

	var rdeps []string
	if transitive {
		rdeps = getDepsTransitive(adjacencyList, targets, depth)
	} else {
		rdeps = getDepsDirect(adjacencyList, targets)
	}

	if reflexive {
		reflexiveTargets := getReflexiveTargets(targets, adjacencyList)
		rdeps = append(rdeps, reflexiveTargets...)
	}
	slices.Sort(rdeps)
	rdeps = slices.Compact(rdeps)
	return rdeps, nil
}
