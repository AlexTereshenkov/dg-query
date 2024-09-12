/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/stretchr/testify/assert"
)

func TestMetricsDependenciesDirect(t *testing.T) {
	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	type TestCase struct {
		input    []byte
		expected map[string]int
	}

	cases := []TestCase{
		// base case
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"bar.py": [
					"eggs.py",
					"baz.py"
				],
				"baz.py": [
					"baz-dep1.py",
					"baz-dep2.py",
					"baz-dep3.py"
				]
			}		
			`),
			expected: map[string]int{
				"foo.py": 1,
				"bar.py": 2,
				"baz.py": 3,
			},
		},
		// empty dependencies for some target
		{
			input: []byte(`
			{
				"foo.py": [],
				"bar.py": ["baz.py"]
			}		
			`),
			expected: map[string]int{
				"foo.py": 0,
				"bar.py": 1,
			},
		},
		// single node
		{
			input: []byte(`
			{
				"foo.py": ["bar.py", "baz.py"]
			}		
			`),
			expected: map[string]int{
				"foo.py": 2,
			},
		},
	}

	for _, testCase := range cases {
		// mocking function that reads a file from disk
		cmd.ReadFile = func(filePath string) []byte {
			return testCase.input
		}
		cmd.RootCmd.SetArgs([]string{"metrics", "--metric=deps-direct", "--dg=dg.json"})
		cmd.RootCmd.Execute()
		var actualOutput map[string]map[string]int
		json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["deps-direct"])
		buf.Reset()
	}
}
