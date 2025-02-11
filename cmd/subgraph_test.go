/*
Copyright © 2025 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseSubgraph struct {
	input    []byte
	rootNode string
	expected string
}

var baseCase string = `
{
	"A": ["B", "B1"], 
	"B": ["C", "C1"],
	"C": ["D"],
	"C1": ["D1"]
}
`

// A <--> C
var baseCaseWithCycle string = `
{
	"A": ["B", "B1"], 
	"B": ["C", "C1"],
	"C": ["A"],
	"C1": ["D1"]
}
`

func TestExtractSubgraph(t *testing.T) {
	cases := []testCaseSubgraph{
		// same roots -> same subgraph
		{
			input:    []byte(baseCase),
			rootNode: "A",
			expected: baseCase,
		},
		// a lower root
		{
			input:    []byte(baseCase),
			rootNode: "B",
			expected: `{
					"B": ["C", "C1"],
					"C": ["D"],
					"C1": ["D1"]
				}`,
		},
		// root is the last node
		{
			input:    []byte(baseCase),
			rootNode: "C1",
			expected: `{
					"C1": ["D1"]
				}`,
		},
		// root is the node without any dependencies -> empty subgraph
		{
			input:    []byte(baseCase),
			rootNode: "D",
			expected: `{}`,
		},
		// graph with a cycle, same root, same subgraph
		{
			input:    []byte(baseCaseWithCycle),
			rootNode: "A",
			expected: baseCaseWithCycle,
		},
		// graph with a cycle, root is at the end of the cycle
		{
			input:    []byte(baseCaseWithCycle),
			rootNode: "C",
			expected: `{
						"A": ["B", "B1"], 
						"B": ["C", "C1"],
						"C": ["A"],
						"C1": ["D1"]
					   }`,
		},
		// graph with a cycle, root is in the middle of the cycle
		{
			input:    []byte(baseCaseWithCycle),
			rootNode: "B",
			expected: `{
						"A": ["B", "B1"], 
						"B": ["C", "C1"],
						"C": ["A"],
						"C1": ["D1"]
					   }`,
		},
		// graph with a cycle, root is outside of the cycle
		{
			input:    []byte(baseCaseWithCycle),
			rootNode: "C1",
			expected: `{"C1": ["D1"]}`,
		},
		// graph with a cycle, root is outside of the cycle
		{
			input: []byte(`{
						"A": ["B", "B1"], 
						"B": ["A", "C", "C1"],
						"C": ["D"],
						"C1": ["D1"],
						"D": ["E"]
						}`),
			rootNode: "C",
			expected: `{"C": ["D"], "D": ["E"]}`,
		},
		// graph is a single cycle, root is either of the nodes
		{
			input:    []byte(`{"A": ["B"], "B": ["A"]}`),
			rootNode: "B",
			expected: `{"A": ["B"], "B": ["A"]}`,
		},
	}
	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		result, err := extractSubgraph("mock-dg.json", testCase.rootNode, MockReadFile)
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
