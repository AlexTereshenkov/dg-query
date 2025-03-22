/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"encoding/json"
	"slices"
	"testing"
	"time"

	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func createAdjacencyLists(nodesCount int) map[string][]string {
	graph := make(map[string][]string)
	for i := 1; i <= nodesCount; i++ {
		node := cast.ToString(i)
		if i < nodesCount {
			nextNode := cast.ToString(i + 1)
			graph[node] = []string{nextNode}
		} else {
			graph[node] = []string{}
		}
	}
	return graph
}

/*
Testing performance of getting dependencies for a node in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}
*/
func TestDependenciesCommandPerfDeepGraph(t *testing.T) {

	startTime := time.Now()
	nodesCount := 10000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	result, err := cmd.Dependencies("mock.json", []string{"1"}, true, false, 0, MockReadFile)
	if err != nil {
		t.Fail()
	}
	expected := make([]string, nodesCount-1)
	for i := range expected {
		expected[i] = cast.ToString(i + 2)
	}

	assert.ElementsMatch(t, expected, result, "Failing assertion")
	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 5 {
		t.Fatalf("Getting dependencies transitively out of a large graph took too long: %s.", elapsedTime)
	}
}

/*
Testing performance of getting dependencies for a node in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}
*/
func TestDependenciesCommandPerfDeepGraphDepthLimit(t *testing.T) {

	startTime := time.Now()
	nodesCount := 1000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	transitive := true
	reflexive := false
	depth := 512
	result, err := cmd.Dependencies("mock.json", []string{"1"}, transitive, reflexive, depth, MockReadFile)
	if err != nil {
		t.Fail()
	}
	expected := make([]string, 512)
	for i := range expected {
		expected[i] = cast.ToString(i + 2)
	}

	assert.ElementsMatch(t, expected, result, "Failing assertion")
	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 5 {
		t.Fatalf("Getting dependencies transitively out of a large graph took too long: %s.", elapsedTime)
	}

}

/*
Testing performance of getting dependents for a node in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}
*/
func TestDependentsCommandPerfDeepGraph(t *testing.T) {
	startTime := time.Now()

	// mocking function that reads a file from disk
	nodesCount := 10000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	transitive := true
	reflexive := false
	depth := 0
	result, err := cmd.Dependents("mock-dg.json", "", []string{"10000"}, transitive, reflexive, depth, MockReadFile)
	if err != nil {
		t.Fail()
	}

	expected := make([]string, nodesCount-1)
	for i := range expected {
		expected[i] = cast.ToString(i + 1)
	}
	slices.Sort(expected)
	assert.ElementsMatch(t, expected, result, "Failing assertion")

	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 2 {
		t.Fatalf("Getting dependents transitively out of a large graph took too long: %s.", elapsedTime)
	}

}

/*
Testing performance of counting dependencies for all nodes in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}.
Despite memoization of intermediate nodes counting, it will
perform poorly on a graph with nodes that have thousands of nested
level dependencies. This is unrealistic in a dependency graph of a
production codebase where the import depth rarely goes over a few
dozen levels until reaching a leaf.
*/
func TestMetricsCommandPerfDeepGraph(t *testing.T) {

	startTime := time.Now()
	nodesCount := 1000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	result, err := cmd.Metrics("mock.json", "", []string{cmd.MetricDependenciesTransitive}, MockReadFile)
	if err != nil {
		t.Fail()
	}
	var actualOutput map[string]map[string]int
	json.Unmarshal(result, &actualOutput)

	assert.Equal(t, 999, actualOutput["deps-transitive"]["1"], "Failing assertion")

	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 5 {
		t.Fatalf("Getting metrics for dependencies count transitively out of a large graph took too long: %s.", elapsedTime)
	}
}

/*
Testing performance of getting paths for a node in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}
*/
func TestPathsCommandPerfDeepGraph(t *testing.T) {

	startTime := time.Now()
	nodesCount := 10000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	result, err := cmd.Paths("mock.json", "1", cast.ToString(nodesCount), 0, MockReadFile)
	if err != nil {
		t.Fail()
	}
	expected := make([]string, nodesCount)
	for i := 0; i < nodesCount; i++ {
		expected[i] = cast.ToString(i + 1)
	}

	assert.ElementsMatch(t, expected, result[0], "Failing assertion")
	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 1 {
		t.Fatalf("Getting paths in a large graph took too long: %s.", elapsedTime)
	}
}

/*
Testing performance of getting subgraph for a root node in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}
*/
func TestSubgraphCommandPerfDeepGraph(t *testing.T) {

	startTime := time.Now()
	nodesCount := 10000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	result, err := cmd.ExtractSubgraph("mock.json", "1", MockReadFile)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 10000, len(result), "Failing assertion")
	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 1 {
		t.Fatalf("Getting subgraph in a large graph took too long: %s.", elapsedTime)
	}
}

/*
Testing performance of getting connected components for a root node in a
deeply nested graph, i.e. {1: [2], 2: [3], 3: [4]..., N: [N+1]}
*/
func TestConnectedComponentsCommandPerfDeepGraph(t *testing.T) {

	startTime := time.Now()
	nodesCount := 10000
	MockReadFile := func(filePath string) ([]byte, error) {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists, nil
	}
	result, err := cmd.ListConnectedComponents("mock.json", MockReadFile)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 10000, len(result[0]), "Failing assertion")
	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 1 {
		t.Fatalf("Getting subgraph in a large graph took too long: %s.", elapsedTime)
	}
}
