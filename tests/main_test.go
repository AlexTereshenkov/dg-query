/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/stretchr/testify/assert"
)

func TestDependenciesDirect(t *testing.T) {

	var buf bytes.Buffer
	// redirection is not required for any subcommands, but this is how it's done for the reference:
	// for _, c := range cmd.RootCmd.Commands() {
	// 	c.SetOut(&buf)
	// }
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cases := []struct {
		input    []byte
		expected []string
		targets  []string
	}{
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
		// mocking function that reads a file from disk
		cmd.ReadFile = func(filePath string) []byte {
			return testCase.input
		}
		cmd.RootCmd.SetArgs(append([]string{"dependencies", "--dg=dg.json"}, testCase.targets...))
		cmd.RootCmd.Execute()
		actualOutput := strings.Split(buf.String(), "\n")[:len(testCase.expected)]
		assert.Equal(t, testCase.expected, actualOutput)
		buf.Reset()
	}
}

func TestDependenciesTransitive(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cases := []struct {
		input    []byte
		expected []string
		targets  []string
	}{
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
			expected: []string{"bar.py"},
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
		// mocking function that reads a file from disk
		cmd.ReadFile = func(filePath string) []byte {
			return testCase.input
		}
		cmd.RootCmd.SetArgs(append([]string{"dependencies", "--transitive", "--dg=dg.json"}, testCase.targets...))
		cmd.RootCmd.Execute()
		actualOutput := strings.Split(buf.String(), "\n")[:len(testCase.expected)]
		assert.Equal(t, testCase.expected, actualOutput)
		buf.Reset()
	}
}
