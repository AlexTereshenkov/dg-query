/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"bytes"
	"encoding/json"
	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
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

	// graphTypes := []func(nodesCount int)map[string][]string{createAdjacencyLists}

	startTime := time.Now()
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	// mocking function that reads a file from disk
	nodesCount := 100000
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
