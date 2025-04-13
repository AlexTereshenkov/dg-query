/*
Copyright © 2025 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseSimplifyTransitiveReduction struct {
	input    []byte
	expected string
}

func TestSimplifyTransitiveReduction(t *testing.T) {
	cases := []testCaseSimplifyTransitiveReduction{
		// empty
		{
			input:    []byte(`{}`),
			expected: `{}`,
		},
		// node with no dependencies
		{
			input:    []byte(`{"A": []}`),
			expected: `{"A": []}`,
		},
		// node with a single dependency with a dependent listed in the adjacency list
		{
			input:    []byte(`{"A": ["B"], "B": []}`),
			expected: `{"A": ["B"], "B": []}`,
		},
		// node with a single dependency with a dependent not listed in the adjacency list
		{
			input:    []byte(`{"A": ["B"]}`),
			expected: `{"A": ["B"]}`,
		},
		// node with dependency on itself
		{
			input:    []byte(`{"A": ["A"]}`),
			expected: `{"A": ["A"]}`,
		},
		// circular dependency
		{
			input:    []byte(`{"A": ["B"], "B": ["A"]}`),
			expected: `{"A": ["B"], "B": ["A"]}`,
		},
		//two  circular dependencies
		{
			input:    []byte(`{"A": ["A"], "B": ["B"]}`),
			expected: `{"A": ["A"], "B": ["B"]}`,
		},
		// base case
		{
			input: []byte(`{
				"A": ["B", "C", "D", "E"], 
				"B": ["D"],
				"C": ["D", "E"],
				"D": ["E"],
				"E": []
			}`),
			expected: `{
					"A": ["B", "C"],
					"B": ["D"],
					"C": ["D"],
					"D": ["E"],
					"E": []
				}`,
		},
		// base case, preserve the order of nodes
		{
			input: []byte(`{
				"B": [
					"C",
					"D"
				],
				"A": [
					"B",
					"D"
				],
				"C": [
					"D"
				]
			}`),
			expected: `{
				"B": [
					"C"
				],
				"A": [
					"B"
				],
				"C": [
					"D"
				]
				}`,
		},
	}
	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		result, err := simplifyAdjacencyList("mock-dg.json", MockReadFile, TechniqueTransitiveReduction)
		if err != nil {
			t.Fail()
		}

		adjacencyListExpected, err := loadJsonFile([]byte(testCase.expected))
		if err != nil {
			t.Fail()
		}
		assert.Equal(t, adjacencyListExpected, result)
	}
}
