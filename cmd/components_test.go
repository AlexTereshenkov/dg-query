/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCaseComponents struct {
	input    []byte
	expected [][]string
}

func TestGetConnectedComponents(t *testing.T) {
	cases := []TestCaseComponents{
		// base case
		{
			input: []byte(`
			{
				"foo": ["bar"],
				"bar": ["spam"],
				"baz": ["eggs"],
				"eggs": ["cheese"]
			}
			`),
			expected: [][]string{
				{"bar", "foo", "spam"},
				{"baz", "cheese", "eggs"},
			},
		},
		// empty graph
		{
			input:    []byte(`{}`),
			expected: [][]string{},
		},
		// one component
		{
			input: []byte(`
			{
				"foo": ["bar"],
				"bar": ["spam"],
				"spam": ["eggs"],
				"eggs": ["cheese"],
				"cheese": []
			}
			`),
			expected: [][]string{{"bar", "cheese", "eggs", "foo", "spam"}},
		},
		// one component
		{
			input: []byte(`
			{
				"foo": []
			}
			`),
			expected: [][]string{{"foo"}},
		},
		// three components with each component being a single node
		{
			input: []byte(`
			{
				"foo": [],
				"bar": [],
				"baz": []
			}
			`),
			expected: [][]string{{"bar"}, {"baz"}, {"foo"}},
		},
		// two components
		{
			input: []byte(`
			{
				"foo": ["bar", "baz"],
				"bar": ["eggs", "spam"],
				"baz": ["ham", "beans"],

				"foo1": ["bar1", "baz1"],
				"bar1": ["eggs1", "spam1"],
				"baz1": ["ham1", "beans1"]
			}
			`),
			expected: [][]string{
				{"bar", "baz", "beans", "eggs", "foo", "ham", "spam"},
				{"bar1", "baz1", "beans1", "eggs1", "foo1", "ham1", "spam1"},
			},
		},
		// three components with components of varying size
		{
			input: []byte(`
			{
				"foo": ["bar", "baz"],
				"bar": ["baz"],
				"ham": ["beans"],
				"cheese": []
			}
			`),
			expected: [][]string{
				{"bar", "baz", "foo"},
				{"beans", "ham"},
				{"cheese"},
			},
		},
		// complex interconnected graph with four components (with one component being a cycle)
		{
			input: []byte(`
			{
				"a": ["b", "c"],
				"b": ["d"],
				"e": ["f"],
				"g": ["h"],
				"h": ["i"],
				"i": ["g"],
				"j": []
			}
			`),
			expected: [][]string{
				{"a", "b", "c", "d"},
				{"e", "f"},
				{"g", "h", "i"},
				{"j"},
			},
		},
		// one cycle is one component
		{
			input: []byte(`
			{
				"foo": ["bar"],
				"bar": ["baz"],
				"baz": ["foo"]
			}
			`),
			expected: [][]string{{"bar", "baz", "foo"}},
		},
		// one cycle with the node connected to itself is still one component
		{
			input: []byte(`
			{
				"foo": ["foo"]
			}
			`),
			expected: [][]string{{"foo"}},
		},
		// two cycles is two components
		{
			input: []byte(`
			{
				"foo": ["bar"],
				"bar": ["foo"],
				"baz": ["spam"],
				"spam": ["baz"]
			}
			`),
			expected: [][]string{
				{"bar", "foo"},
				{"baz", "spam"},
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		components, err := listConnectedComponents("mock.json", MockReadFile)
		if err != nil {
			t.Fail()
		}
		assert.Equal(t, testCase.expected, components)
	}
}
