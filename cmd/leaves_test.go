/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseLeaves struct {
	input    []byte
	expected []string
}

func TestLeaves(t *testing.T) {
	cases := []testCaseLeaves{
		// empty graph
		{
			input:    []byte(`{}`),
			expected: []string{},
		},
		// node with no dependencies
		{
			input: []byte(`
		{
			"foo.py": []		
		}		
		`),
			expected: []string{"foo.py"},
		},
		// node with dependency on itself
		{
			input: []byte(`
		{
			"foo.py": ["foo.py"]		
		}		
		`),
			expected: []string{},
		},
		// circular dependency
		{
			input: []byte(`
		{
			"foo.py": ["bar.py"],
			"bar.py": ["foo.py"]		
		}		
		`),
			expected: []string{},
		},
		//two  circular dependencies
		{
			input: []byte(`
		{
			"foo.py": ["foo.py"],
			"bar.py": ["bar.py"]		
		}		
		`),
			expected: []string{},
		},
		// nodes with some dependencies
		{
			input: []byte(`
		{
			"foo.py": [
				"bar.py",
				"baz.py"
			],
			"spam.py": [
				"eggs.py",
				"ham.py"
			]		
		}		
		`),
			expected: []string{"foo.py", "spam.py"},
		},
		// linked list dependencies
		{
			input: []byte(`
		{
				"1": ["2"],
				"2": ["3"],
				"3": ["4"],
				"4": ["5"]		
		}		
		`),
			expected: []string{"1"},
		},
		// complex graph with some nodes that are not in the keys of the adjacency list;
		// note "d", "e", and "g" that some nodes depend on, but they are not declared
		// ensure the leaves are returned in sorted order
		{
			input: []byte(`
		{
				"f": ["g"],
				"c": ["d", "e"],
				"h": [],
				"a": ["b", "c"],
				"b": ["d"]				
		}		
		`),
			expected: []string{"a", "f", "h"},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		result, err := leaves("mock.json", MockReadFile)
		if err != nil {
			t.Fail()
		}
		assert.Equal(t, testCase.expected, result)
	}
}
