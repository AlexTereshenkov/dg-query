/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCasePaths struct {
	input      []byte
	expected   [][]string
	fromTarget string
	toTarget   string
}

func TestPaths(t *testing.T) {

	cases := []testCasePaths{
		/*

				    	A
			           / \
			          B   B1
			          |    |
			          C   C1
			           \  /
			            D

		*/
		{
			input: []byte(`{
				"A": ["B", "B1"], 
				"B": ["C"], 
				"B1": ["C1"], 
				"C": ["D"],
				"C1": ["D"]
			}`),
			expected: [][]string{
				{"A", "B", "C", "D"},
				{"A", "B1", "C1", "D"},
			},
			fromTarget: "A",
			toTarget:   "D",
		},
		/*
		       A
		      / \
		     B   B1
		    / \ / \
		   C---D---C1
		*/
		{
			input: []byte(`{
				"A": ["B", "B1"], 
				"B": ["C", "D"], 
				"B1": ["C1", "D"], 
				"C": ["D"],
				"C1": ["D"]
			}`),
			expected: [][]string{
				{"A", "B", "D"},
				{"A", "B1", "D"},
				{"A", "B", "C", "D"},
				{"A", "B1", "C1", "D"},
			},
			fromTarget: "A",
			toTarget:   "D",
		},
		/*
		       A
		      / \
		     B   B1
		    / \ / \
		   C---D---C1
		       |
		       E
		*/
		{
			input: []byte(`{
				"A": ["B", "B1"], 
				"B": ["C", "D"], 
				"B1": ["C1", "D"], 
				"C": ["D"],
				"C1": ["D"],
				"D": ["E"]
			}`),
			expected: [][]string{
				{"B", "D", "E"},
				{"B", "C", "D", "E"},
			},
			fromTarget: "B",
			toTarget:   "E",
		},
		/*
		       A
		      / \
		     B   B1
		    / \ / \
		   C---D---C1
		       |    \
		       E     F
		*/
		{
			input: []byte(`{
				"A": ["B", "B1"], 
				"B": ["C", "D"], 
				"B1": ["C1", "D"], 
				"C": ["D"],
				"C1": ["D", "F"],
				"D": ["E"]
			}`),
			// there's no path between these two targets
			expected:   [][]string{},
			fromTarget: "B",
			toTarget:   "F",
		},
		/*
			A -> B -> C -> D -|
			^				  |
			|                 |
			-------------------
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["D"], 
				"D": ["A"]
			}`),
			expected: [][]string{
				{"A", "B", "C", "D"},
			},
			fromTarget: "A",
			toTarget:   "D",
		},
		/*
			A -> B -> C -> D -|
			^				  |
			|                 |
			-------------------
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["D"], 
				"D": ["A"]
			}`),
			expected: [][]string{
				{"B", "C", "D", "A"},
			},
			fromTarget: "B",
			toTarget:   "A",
		},
		/*
			A -> B -> C -> D -|- E
			^				  |
			|                 |
			-------------------
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["D"], 
				"D": ["A", "E"]
			}`),
			expected: [][]string{
				{"A", "B", "C", "D", "E"},
			},
			fromTarget: "A",
			toTarget:   "E",
		},
		/*
			A -> B -|
			^		|
			|       |
			--------
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["A"]
			}`),
			expected: [][]string{
				{"A", "B"},
			},
			fromTarget: "A",
			toTarget:   "B",
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		result, err := paths("mock-dg.json", testCase.fromTarget, testCase.toTarget, MockReadFile)
		if err != nil {
			t.Fail()
		}
		// the order of paths may change depending on the implementation of the DFS,
		// but the order the paths are returned in shouldn't really matter
		assert.ElementsMatch(t, testCase.expected, result)
	}
}
