/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"bytes"
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/stretchr/testify/assert"
)

func createAdjacencyLists(nodesCount int) map[string][]string {
	graph := make(map[string][]string)
	for i := 1; i <= nodesCount; i++ {
		node := strconv.Itoa(i)
		if i < nodesCount {
			nextNode := strconv.Itoa(i + 1)
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
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	// mocking function that reads a file from disk
	nodesCount := 10000
	cmd.ReadFile = func(filePath string) []byte {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists
	}
	cmd.RootCmd.SetArgs([]string{"dependencies", "--transitive", "--dg=dg.json", "1"})
	cmd.RootCmd.Execute()

	expected := make([]string, nodesCount-1)
	for i := range expected {
		expected[i] = strconv.Itoa(i + 2)
	}

	actualOutput := strings.Split(buf.String(), "\n")[:nodesCount-1]
	assert.ElementsMatch(t, expected, actualOutput, "Failing assertion")
	buf.Reset()

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
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	// mocking function that reads a file from disk
	nodesCount := 10000
	cmd.ReadFile = func(filePath string) []byte {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists
	}
	cmd.RootCmd.SetArgs([]string{"dependents", "--transitive", "--dg=dg.json", "10000"})
	cmd.RootCmd.Execute()

	expected := make([]string, nodesCount-1)
	for i := range expected {
		expected[i] = strconv.Itoa(i + 1)
	}
	slices.Sort(expected)

	actualOutput := strings.Split(buf.String(), "\n")[:nodesCount-1]
	assert.ElementsMatch(t, expected, actualOutput, "Failing assertion")
	buf.Reset()

	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 5 {
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
dozen levels.
*/
func TestMetricsCommandPerfDeepGraph(t *testing.T) {

	startTime := time.Now()
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	// mocking function that reads a file from disk
	nodesCount := 1000
	cmd.ReadFile = func(filePath string) []byte {
		lists, _ := json.Marshal(createAdjacencyLists(nodesCount))
		return lists
	}
	cmd.RootCmd.SetArgs([]string{"metrics", "--metric=deps-transitive", "--dg=dg.json"})
	cmd.RootCmd.Execute()

	var actualOutput map[string]map[string]int
	json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &actualOutput)

	assert.Equal(t, 999, actualOutput["deps-transitive"]["1"], "Failing assertion")
	buf.Reset()

	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() > 5 {
		t.Fatalf("Getting metrics for dependencies count transitively out of a large graph took too long: %s.", elapsedTime)
	}
}
