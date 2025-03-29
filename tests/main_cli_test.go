package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/stretchr/testify/assert"
)

func TestCliDependencies(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs(append([]string{"dependencies", "--transitive", "--reflexive", "--dg=examples/dg.json"}, "foo.py"))
	cmd.RootCmd.Execute()

	expected := []string{"foo-dep1-dep1.py", "foo-dep1-dep2.py", "foo-dep1.py", "foo-dep2.py", "foo.py"}
	actualOutput := strings.Split(buf.String(), "\n")[:len(expected)]
	assert.Equal(t, expected, actualOutput)
	buf.Reset()
}

func TestCliDependents(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs(append([]string{"dependents", "--transitive", "--dg=examples/dg.json"}, "foo-dep1-dep1.py"))
	cmd.RootCmd.Execute()

	expected := []string{"foo-dep1.py", "foo.py"}
	actualOutput := strings.Split(buf.String(), "\n")[:len(expected)]
	assert.Equal(t, expected, actualOutput)
	buf.Reset()
}

func TestCliLeaves(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs([]string{"leaves", "--dg=examples/dg.json"})
	cmd.RootCmd.Execute()

	expected := []string{"foo.py", "spam.py"}
	actualOutput := strings.Split(buf.String(), "\n")[:len(expected)]
	assert.Equal(t, expected, actualOutput)
	buf.Reset()
}

func TestCliPaths(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs([]string{"paths", "--dg=examples/dg.json", "--from=foo.py", "--to=foo-dep1-dep1.py", "--n=1"})
	cmd.RootCmd.Execute()

	expected := []byte(`[["foo.py", "foo-dep1.py", "foo-dep1-dep1.py"]]`)

	var actualOutput [][]string
	var expectedOutput [][]string
	json.Unmarshal(buf.Bytes(), &actualOutput)
	json.Unmarshal(expected, &expectedOutput)
	assert.Equal(t, expectedOutput, actualOutput)
	buf.Reset()
}

func TestCliCycles(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs([]string{"cycles", "--dg=examples/dg.json"})
	cmd.RootCmd.Execute()

	expected := []byte(`[]`)

	var actualOutput [][]string
	var expectedOutput [][]string
	json.Unmarshal(buf.Bytes(), &actualOutput)
	json.Unmarshal(expected, &expectedOutput)
	assert.Equal(t, expectedOutput, actualOutput)
	buf.Reset()
}

func TestCliSubgraph(t *testing.T) {
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs([]string{"subgraph", "--dg=examples/dg.json", "--root=foo-dep1.py"})
	cmd.RootCmd.Execute()

	expected := []byte(`{"foo-dep1.py": ["foo-dep1-dep1.py","foo-dep1-dep2.py"]}`)
	var actualOutput cmd.AdjacencyList
	var expectedOutput cmd.AdjacencyList
	json.Unmarshal(buf.Bytes(), &actualOutput)
	json.Unmarshal(expected, &expectedOutput)
	assert.Equal(t, expectedOutput, actualOutput)
	buf.Reset()
}

func TestCliComponents(t *testing.T) {
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs([]string{"components", "--dg=examples/dg.json"})
	cmd.RootCmd.Execute()

	expected := []byte(`
	[
		["foo-dep1-dep1.py","foo-dep1-dep2.py","foo-dep1.py","foo-dep2.py","foo.py"],
		["spam-dep1.py","spam-dep2-dep1.py","spam-dep2-dep2.py","spam-dep2.py","spam.py"]
	]
	`)

	var actualOutput [][]string
	var expectedOutput [][]string
	json.Unmarshal(buf.Bytes(), &actualOutput)
	json.Unmarshal(expected, &expectedOutput)
	assert.Equal(t, expectedOutput, actualOutput)
	buf.Reset()
}

func TestCliMetrics(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs([]string{"metrics",
		"--metric=deps-direct", "--metric=deps-transitive", "--metric=rdeps-direct", "--metric=rdeps-transitive", "--metric=components-count",
		"--dg=examples/dg.json"})
	cmd.RootCmd.Execute()

	expected := []byte(`
		{
		"components-count": {
			"count": 2
		},
		"deps-direct": {
			"foo-dep1.py": 2,
			"foo.py": 2,
			"spam-dep2.py": 2,
			"spam.py": 2,
			"foo-dep2.py": 0,
			"spam-dep1.py": 0,
			"foo-dep1-dep1.py": 0,
			"foo-dep1-dep2.py": 0,
			"spam-dep2-dep1.py": 0,
			"spam-dep2-dep2.py": 0
		},
		"deps-transitive": {
			"foo-dep1.py": 2,
			"foo.py": 4,
			"spam-dep2.py": 2,
			"spam.py": 4,
			"foo-dep2.py": 0,
			"spam-dep1.py": 0,
			"foo-dep1-dep1.py": 0,
			"foo-dep1-dep2.py": 0,
			"spam-dep2-dep1.py": 0,
			"spam-dep2-dep2.py": 0
		},
		"rdeps-direct": {
			"foo-dep1-dep1.py": 1, 
			"foo-dep1-dep2.py": 1, 
			"foo-dep1.py": 1, 
			"foo-dep2.py": 1, 
			"foo.py": 0, 
			"spam-dep1.py": 1, 
			"spam-dep2-dep1.py": 1, 
			"spam-dep2-dep2.py": 1, 
			"spam-dep2.py": 1, 
			"spam.py": 0
		}, 
		"rdeps-transitive": {
			"foo-dep1-dep1.py": 2, 
			"foo-dep1-dep2.py": 2, 
			"foo-dep1.py": 1, 
			"foo-dep2.py": 1, 
			"foo.py": 0, 
			"spam-dep1.py": 1, 
			"spam-dep2-dep1.py": 2, 
			"spam-dep2-dep2.py": 2, 
			"spam-dep2.py": 1, 
			"spam.py": 0
		}
		}
	`)
	var actualOutput map[string]map[string]int
	var expectedOutput map[string]map[string]int
	json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &actualOutput)
	json.Unmarshal(expected, &expectedOutput)

	assert.Equal(t, expectedOutput, actualOutput)
	buf.Reset()
}
