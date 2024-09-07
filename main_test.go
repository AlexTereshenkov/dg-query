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
		// plain
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
		assert.ElementsMatch(t, testCase.expected, actualOutput, "Failing assertion")
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
		// plain transitive dependencies
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
			expected: []string{"bar.py", "baz.py", "eggs.py", "ham.py", "cheese.py"},
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
		// only circular dependencies (both circle targets)
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
			targets:  []string{"foo.py", "bar,py"},
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
		assert.ElementsMatch(t, testCase.expected, actualOutput, "Failing assertion")
		buf.Reset()
	}
}
