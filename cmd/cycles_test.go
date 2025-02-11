/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseCycles struct {
	input    []byte
	expected [][]string
}

func TestFindCycles(t *testing.T) {

	cases := []testCaseCycles{
		// no cycles
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"]
			}`),
			expected: [][]string{},
		},
		// one cycle with no other nodes
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["A"]
			}`),
			expected: [][]string{
				{"A", "B"},
			},
		},
		// one cycle with the node connected to itself
		{
			input: []byte(`
            {
                "A": ["A"]
            }
            `),
			expected: [][]string{
				{"A"},
			},
		},
		// two cycles with the nodes connected to themselves
		{
			input: []byte(`
            {
                "A": ["A"],
				"B": ["B"]
            }
            `),
			expected: [][]string{
				{"A"}, {"B"},
			},
		},
		// one cycle among nodes outside of the cycle
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["A"],
                "D": ["D1", "D2"]
			}`),
			expected: [][]string{
				{"A", "B", "C"},
			},
		},
		// two cycles
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["A"], 
				"C": ["D"],
                "D": ["C"]
			}`),
			expected: [][]string{
				{"A", "B"}, {"C", "D"},
			},
		},
		/* two cycles with a node in between
		A -> B -> C  -> D -> E -> F -> G
		^	 	  |   		 ^  	   |
		|         |          |         |
		-----<-----          -----<-----
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["A", "D"],
                "D": ["E"],
				"E": ["F"],
				"F": ["G"],
				"G": ["E"]
			}`),
			expected: [][]string{
				{"A", "B", "C"}, {"E", "F", "G"},
			},
		},
		/* intervened cycles
		A -> B -> C  -> D
		^	 ^	  |     |
		|    |    |     |
		-----------------
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["A", "D"],
                "D": ["B"]
			}`),
			expected: [][]string{
				{"A", "B", "C"}, {"B", "C", "D"},
			},
		},
		/* inner cycle overlapping with an outer cycle (shared start node)
		A -> B -> C  -> D -> E
		^	 	  |     	 |
		|         |          |
		----<----------<-----
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["A", "D"],
                "D": ["E"],
				"E": ["A"]
			}`),
			expected: [][]string{
				{"A", "B", "C"}, {"A", "B", "C", "D", "E"},
			},
		},
		/* inner cycle inside an outer cycle (distinct start nodes)
			A -> B -> C  -> D -> E
			^	 ^	 	    |    |
		    |	 |          |    |
			----<----------<------
		*/
		{
			input: []byte(`{
				"A": ["B"], 
				"B": ["C"], 
				"C": ["D"],
                "D": ["B", "E"],
				"E": ["A"]
			}`),
			expected: [][]string{
				{"B", "C", "D"}, {"A", "B", "C", "D", "E"},
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		result, err := cycles("mock-dg.json", MockReadFile)
		if err != nil {
			t.Fail()
		}
		// the order of cycles may change depending on the implementation of the DFS,
		// but the order the cycles are returned in shouldn't really matter
		assert.ElementsMatch(t, testCase.expected, result)
	}
}
