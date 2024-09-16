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

func TestDependentsDirect(t *testing.T) {

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
		// mocking function that reads a file from disk
		cmd.ReadFile = func(filePath string) []byte {
			return testCase.input
		}
		cmd.RootCmd.SetArgs(append([]string{"dependents", "--dg=dg.json"}, testCase.targets...))
		cmd.RootCmd.Execute()
		actualOutput := strings.Split(buf.String(), "\n")[:len(testCase.expected)]
		assert.Equal(t, testCase.expected, actualOutput)
		buf.Reset()
	}
}

func TestDependentsTransitive(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cases := []struct {
		input    []byte
		expected []string
		targets  []string
	}{
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
			expected: []string{"bar.py", "baz.py", "foo.py"},
			targets:  []string{"spam.py"},
		},
	}

	for _, testCase := range cases {
		// mocking function that reads a file from disk
		cmd.ReadFile = func(filePath string) []byte {
			return testCase.input
		}
		cmd.RootCmd.SetArgs(append([]string{"dependents", "--transitive", "--dg=dg.json"}, testCase.targets...))
		cmd.RootCmd.Execute()
		actualOutput := strings.Split(buf.String(), "\n")[:len(testCase.expected)]
		assert.Equal(t, testCase.expected, actualOutput)
		buf.Reset()
	}
}

func TestDependentsTransitiveReflexiveClosure(t *testing.T) {

	var buf bytes.Buffer
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	cases := []struct {
		input    []byte
		expected []string
		targets  []string
	}{
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
		// mocking function that reads a file from disk
		cmd.ReadFile = func(filePath string) []byte {
			return testCase.input
		}
		cmd.RootCmd.SetArgs(append([]string{"dependents", "--transitive", "--reflexive", "--dg=dg.json"}, testCase.targets...))
		cmd.RootCmd.Execute()
		actualOutput := strings.Split(buf.String(), "\n")[:len(testCase.expected)]
		assert.Equal(t, testCase.expected, actualOutput)
		buf.Reset()
	}
}
