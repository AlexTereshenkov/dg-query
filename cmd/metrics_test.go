/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"

	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	input    []byte
	expected map[string]int
}

func TestMetricsDependenciesDirect(t *testing.T) {
	cases := []TestCase{
		// base case
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"bar.py": [
					"eggs.py",
					"baz.py"
				],
				"baz.py": [
					"baz-dep1.py",
					"baz-dep2.py",
					"baz-dep3.py"
				]
			}		
			`),
			expected: map[string]int{
				"foo.py": 1,
				"bar.py": 2,
				"baz.py": 3,
			},
		},
		// empty dependencies for some target
		{
			input: []byte(`
			{
				"foo.py": [],
				"bar.py": ["baz.py"]
			}		
			`),
			expected: map[string]int{
				"foo.py": 0,
				"bar.py": 1,
			},
		},
		// single node
		{
			input: []byte(`
			{
				"foo.py": ["bar.py", "baz.py"]
			}		
			`),
			expected: map[string]int{
				"foo.py": 2,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		metricsItems := []string{MetricDependenciesDirect}
		result := metrics("mock.json", "", metricsItems, MockReadFile)
		var actualOutput map[string]map[string]int
		json.Unmarshal(result, &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["deps-direct"])
	}
}

func TestMetricsDependenciesTransitive(t *testing.T) {
	cases := []TestCase{
		// base case
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"bar.py": [
					"eggs.py",
					"baz.py"
				],
				"baz.py": [
					"baz-dep1.py",
					"baz-dep2.py",
					"baz-dep3.py"
				]
			}
			`),
			expected: map[string]int{
				"foo.py": 1,
				"bar.py": 5,
				"baz.py": 3,
			},
		},
		// cyclic dependencies
		{
			input: []byte(`
			{
				"foo.py": [],
				"bar.py": ["baz.py"],
				"baz.py": ["foo.py"]
			}
			`),
			expected: map[string]int{
				"foo.py": 0,
				"bar.py": 2,
				"baz.py": 1,
			},
		},
		// transitive chain
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"spam.py": [
					"eggs.py"
				],
				"eggs.py": [
					"baz.py"
				]
			}
			`),
			expected: map[string]int{
				"foo.py":  3,
				"spam.py": 2,
				"eggs.py": 1,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		metricsItems := []string{MetricDependenciesTransitive}
		result := metrics("mock.json", "", metricsItems, MockReadFile)
		var actualOutput map[string]map[string]int
		json.Unmarshal(result, &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["deps-transitive"])
	}
}

func TestMetricsReverseDependenciesDirect(t *testing.T) {
	cases := []TestCase{
		// base case
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"bar.py": [
					"spam.py",
					"baz.py"
				],
				"baz.py": [
					"baz-dep1.py"
				]
			}		
			`),
			expected: map[string]int{
				"spam.py":     2,
				"baz.py":      1,
				"baz-dep1.py": 1,
			},
		},
		// empty dependencies for some target
		{
			input: []byte(`
			{
				"foo.py": [],
				"bar.py": ["baz.py"]
			}		
			`),
			expected: map[string]int{
				"baz.py": 1,
			},
		},
		// single node
		{
			input: []byte(`
			{
				"foo.py": ["bar.py", "baz.py"]
			}		
			`),
			expected: map[string]int{
				"bar.py": 1,
				"baz.py": 1,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		metricsItems := []string{MetricReverseDependenciesDirect}
		result := metrics("mock.json", "", metricsItems, MockReadFile)
		var actualOutput map[string]map[string]int
		json.Unmarshal(result, &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["rdeps-direct"])
	}
}

func TestMetricsReverseDependenciesTransitive(t *testing.T) {
	cases := []TestCase{
		// base case
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"bar.py": [
					"spam.py",
					"baz.py"
				],
				"baz.py": [
					"baz-dep1.py"
				]
			}		
			`),
			expected: map[string]int{
				"spam.py":     2,
				"baz.py":      1,
				"baz-dep1.py": 2,
			},
		},
		// cyclic dependencies
		{
			input: []byte(`
			{
				"foo.py": [],
				"bar.py": ["baz.py"],
				"baz.py": ["foo.py"]
			}
			`),
			expected: map[string]int{
				"foo.py": 2,
				"baz.py": 1,
			},
		},
		// transitive chain
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py"
				],
				"spam.py": [
					"eggs.py"
				],
				"eggs.py": [
					"baz.py"
				]
			}
			`),
			expected: map[string]int{
				"baz.py":  3,
				"spam.py": 1,
				"eggs.py": 2,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) []byte {
			return testCase.input
		}
		metricsItems := []string{MetricReverseDependenciesTransitive}
		result := metrics("mock.json", "", metricsItems, MockReadFile)
		var actualOutput map[string]map[string]int
		json.Unmarshal(result, &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["rdeps-transitive"])
	}
}

func TestMetricsCombined(t *testing.T) {
	input := []byte(`
	{
		"foo.py": [
			"spam.py"
		],
		"bar.py": [
			"eggs.py",
			"baz.py"
		],
		"baz.py": [
			"baz-dep1.py",
			"baz-dep2.py",
			"baz-dep3.py"
		]
	}
	`)

	MockReadFile := func(filePath string) []byte {
		return input
	}

	metricsItems := []string{MetricDependenciesDirect, MetricDependenciesTransitive, MetricReverseDependenciesDirect, MetricReverseDependenciesTransitive}
	result := metrics("mock.json", "", metricsItems, MockReadFile)
	var actualOutput map[string]map[string]int
	json.Unmarshal(result, &actualOutput)

	// check that all metrics are present in the report
	for _, metric := range metricsItems {
		_, exists := actualOutput[metric]
		assert.True(t, exists, "Expected metric '%s' to exist in the report", metric)
	}
}
