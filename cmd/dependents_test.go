/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseAdjacencyLists(t *testing.T) {
	dg := AdjacencyList{"foo": {"bar", "baz"}, "spam": {"eggs", "bar"}}
	rdg := reverseAdjacencyLists(dg)
	expected := AdjacencyList{"bar": {"foo", "spam"}, "baz": {"foo"}, "eggs": {"spam"}}
	assert.True(t, reflect.DeepEqual(rdg, expected))
}

type testCaseDependents struct {
	input    []byte
	expected []string
	targets  []string
	// by default zero
	depth int
}

func TestDependentsDirect(t *testing.T) {

	cases := []testCaseDependents{
		// empty graph
		{
			input:    []byte(`{}`),
			expected: []string{},
			targets:  []string{"foo.py"},
		},
		// node with no dependents
		{
			input: []byte(`
		{
			"foo.py": []
		
		}		
		`),
			expected: []string{},
			targets:  []string{"foo.py"},
		},
		// node that is a dependency for itself
		{
			input: []byte(`
		{
			"foo.py": ["foo.py"]
		
		}		
		`),
			expected: []string{"foo.py"},
			targets:  []string{"foo.py"},
		},
		// circular dependency
		{
			input: []byte(`
		{
			"foo.py": ["bar.py"],
			"bar.py": ["foo.py"]
		
		}		
		`),
			expected: []string{"bar.py"},
			targets:  []string{"foo.py"},
		},
		// node with some dependents
		{
			input: []byte(`
			{
				"foo.py": [
					"bar.py",
					"baz.py"
				],
				"spam.py": [
					"ham.py",
					"eggs.py",        
					"bar.py"
				]
			}		
		`),
			expected: []string{"foo.py", "spam.py"},
			targets:  []string{"bar.py"},
		},
		// only some node being a dependent
		{
			input: []byte(`
		{
			"foo.py": [
				"bar.py"
			],
			"baz.py": [
				"spam.py"
			]		
		}		
		`),
			expected: []string{"foo.py"},
			targets:  []string{"bar.py"},
		},
		// asking for non-existing node
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
			expected: []string{},
			targets:  []string{"non-existing.py"},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		transitive := false
		reflexive := false
		// use dependency graph to be reversed
		result := dependents("mock-dg.json", "", testCase.targets, transitive, reflexive, testCase.depth, MockReadFile)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestDependentsTransitiveReflexiveClosure(t *testing.T) {
	cases := []testCaseDependents{
		// node with some dependents
		{
			input: []byte(`
			{
				"foo.py": [
					"baz.py"
				],
				"baz.py": [
					"bar.py"
				],
				"bar.py": [	
					"spam.py"
				]
			}		
		`),
			expected: []string{"bar.py", "baz.py", "foo.py", "spam.py"},
			targets:  []string{"spam.py"},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		transitive := true
		reflexive := true
		result := dependents("mock-dg.json", "", testCase.targets, transitive, reflexive, testCase.depth, MockReadFile)
		assert.Equal(t, testCase.expected, result)
	}
}
