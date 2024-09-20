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

func TestCliMetrics(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cmd.RootCmd.SetArgs(append([]string{"metrics", "--metric=deps-direct", "--metric=deps-transitive", "--dg=examples/dg.json"}, "foo.py"))
	cmd.RootCmd.Execute()

	expected := []byte(`
		{
		"deps-direct": {
			"foo-dep1.py": 2,
			"foo.py": 2,
			"spam-dep2.py": 2,
			"spam.py": 2
		},
		"deps-transitive": {
			"foo-dep1.py": 2,
			"foo.py": 4,
			"spam-dep2.py": 2,
			"spam.py": 4
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
