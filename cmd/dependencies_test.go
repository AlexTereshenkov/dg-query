/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseDependencies struct {
	input    []byte
	expected []string
	targets  []string
}

func TestDependenciesDirect(t *testing.T) {
	cases := []testCaseDependencies{
		// empty graph
		{
			input:    []byte(`{}`),
			expected: []string{},
			targets:  []string{"foo.py"},
		},
		// node with no dependencies
		{
			input: []byte(`
		{
			"foo.py": []
		
		}		
		`),
			expected: []string{},
			targets:  []string{"foo.py"},
		},
		// node with dependency on itself
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
		// node with some dependencies
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
			expected: []string{"bar.py", "baz.py"},
			targets:  []string{"foo.py"},
		},
		// nodes with same dependencies (no duplicates in the output)
		{
			input: []byte(`
		{
			"foo.py": [
				"bar.py"
			],
			"baz.py": [
				"bar.py"
			]		
		}		
		`),
			expected: []string{"bar.py"},
			targets:  []string{"foo.py", "baz.py"},
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
		// asking for multiple nodes
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
			expected: []string{"bar.py", "baz.py", "eggs.py", "ham.py"},
			targets:  []string{"foo.py", "spam.py"},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		transitive := false
		reflexive := false
		result := dependencies("mock.json", testCase.targets, transitive, reflexive, MockReadFile)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestDependenciesTransitive(t *testing.T) {
	cases := []testCaseDependencies{
		// plain transitive dependencies, one level only
		{
			input: []byte(`
			{
				"foo.py": [
					"bar.py",
					"baz.py"
				],
				"bar.py": [],
				"baz.py": []
			}
			`),
			expected: []string{"bar.py", "baz.py"},
			targets:  []string{"foo.py"},
		},
		// plain transitive dependencies, many levels
		{
			input: []byte(`
			{
				"foo.py": [
					"bar.py",
					"baz.py"
				],
				"bar.py": [
					"eggs.py",
					"ham.py"
				],
				"eggs.py": [
					"cheese.py"
				]
			}
			`),
			expected: []string{"bar.py", "baz.py", "cheese.py", "eggs.py", "ham.py"},
			targets:  []string{"foo.py"},
		},
		// nodes with same dependencies (no duplicates in the output)
		{
			input: []byte(`
		{
			"foo.py": [
				"bar.py",
				"spam.py"
			],
			"bar.py": [
				"baz.py",
				"spam.py"
			],
			"spam.py": [
				"baz.py"
			]
		}
		`),
			expected: []string{"bar.py", "baz.py", "spam.py"},
			targets:  []string{"foo.py"},
		},
		// some circular dependencies
		{
			input: []byte(`
			{
				"foo.py": [
					"bar.py",
					"baz.py"
				],
				"bar.py": [
					"foo.py",
					"ham.py"
				],
				"eggs.py": [
					"cheese.py"
				]
			}
			`),
			expected: []string{"bar.py", "baz.py", "foo.py", "ham.py"},
			targets:  []string{"foo.py"},
		},
		// only circular dependencies (single circle target)
		{
			input: []byte(`
			{
				"foo.py": [
					"bar.py"
				],
				"bar.py": [
					"foo.py"
				]
			}
			`),
			expected: []string{"bar.py", "foo.py"},
			targets:  []string{"foo.py"},
		},
		// transitive circular dependency
		{
			input: []byte(`
			{
				"foo.py": [
					"bar.py"
				],
				"bar.py": [
					"baz.py"
				],
				"baz.py": [
					"foo.py"
				]

			}
			`),
			expected: []string{"bar.py", "baz.py", "foo.py"},
			targets:  []string{"foo.py"},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		transitive := true
		reflexive := false
		result := dependencies("mock.json", testCase.targets, transitive, reflexive, MockReadFile)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestDependenciesReflexiveClosure(t *testing.T) {

	cases := []testCaseDependencies{
		// base case
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"bar.py": [
					"eggs.py"
				]
			}
			`),
			expected: []string{"bar.py", "eggs.py", "foo.py", "spam.py"},
			targets:  []string{"foo.py", "bar.py"},
		},
		// empty dependencies with a non-existing target (case 1)
		{
			input: []byte(`
			{
				"foo.py": []
			}
			`),
			expected: []string{"foo.py"},
			targets:  []string{"foo.py", "bar.py"},
		},
		// empty dependencies with a non-existing target (case 2)
		{
			input: []byte(`
			{
				"foo.py": []
			}
			`),
			expected: []string{},
			targets:  []string{"bar.py"},
		},
		// duplicate input targets
		{
			input: []byte(`
			{
				"foo.py": ["bar.py"]
			}
			`),
			expected: []string{"bar.py", "foo.py"},
			targets:  []string{"foo.py", "foo.py"},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		transitive := false
		reflexive := true
		result := dependencies("mock.json", testCase.targets, transitive, reflexive, MockReadFile)
		assert.Equal(t, testCase.expected, result)
	}
}