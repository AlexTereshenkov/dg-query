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
				"foo.py":      1,
				"bar.py":      2,
				"baz.py":      3,
				"spam.py":     0,
				"eggs.py":     0,
				"baz-dep1.py": 0,
				"baz-dep2.py": 0,
				"baz-dep3.py": 0,
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
				"baz.py": 0,
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
				"bar.py": 0,
				"baz.py": 0,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		metricsItems := []string{MetricDependenciesDirect}
		result, err := metrics("mock.json", "", metricsItems, MockReadFile)
		if err != nil {
			t.Fail()
		}
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
				"foo.py":      1,
				"bar.py":      5,
				"baz.py":      3,
				"spam.py":     0,
				"eggs.py":     0,
				"baz-dep1.py": 0,
				"baz-dep2.py": 0,
				"baz-dep3.py": 0,
			},
		},
		// looks like cyclic dependencies
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
		// true like cyclic dependencies
		{
			input: []byte(`
			{
				"foo.py": ["bar.py"],
				"bar.py": ["baz.py"],
				"baz.py": ["foo.py"]
			}
			`),
			expected: map[string]int{
				"foo.py": 3,
				"bar.py": 3,
				"baz.py": 3,
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
				"baz.py":  0,
			},
		},
		// duplicate dependencies should be discarded when counting
		// ("foo.py" and all of its transitive dependencies depend on "spam.py")
		{
			input: []byte(`
			{
				"foo.py": [
					"spam.py",
					"bar.py"
				],
				"bar.py": [
					"eggs.py",
					"spam.py"
				],
				"eggs.py": [
					"baz.py",
					"spam.py"
				]
			}
			`),
			expected: map[string]int{
				"foo.py":  4,
				"bar.py":  3,
				"eggs.py": 2,
				"baz.py":  0,
				"spam.py": 0,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		metricsItems := []string{MetricDependenciesTransitive}
		result, err := metrics("mock.json", "", metricsItems, MockReadFile)
		if err != nil {
			t.Fail()
		}
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
				"foo.py":      0,
				"bar.py":      0,
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
				"bar.py": 0,
				"foo.py": 0,
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
				"foo.py": 0,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		metricsItems := []string{MetricReverseDependenciesDirect}
		result, err := metrics("mock.json", "", metricsItems, MockReadFile)
		if err != nil {
			t.Fail()
		}
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
				"foo.py":      0,
				"bar.py":      0,
			},
		},
		// looks like cyclic dependencies
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
				"baz.py": 1,
				"bar.py": 0,
			},
		},
		// true cyclic dependencies
		{
			input: []byte(`
			{
				"foo.py": ["bar.py"],
				"bar.py": ["baz.py"],
				"baz.py": ["foo.py"]
			}
			`),
			expected: map[string]int{
				"foo.py": 3,
				"baz.py": 3,
				"bar.py": 3,
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
				"eggs.py": 2,
				"spam.py": 1,
				"foo.py":  0,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		metricsItems := []string{MetricReverseDependenciesTransitive}
		result, err := metrics("mock.json", "", metricsItems, MockReadFile)
		if err != nil {
			t.Fail()
		}
		var actualOutput map[string]map[string]int
		json.Unmarshal(result, &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["rdeps-transitive"])
	}
}

func TestMetricsConnectedComponentsCount(t *testing.T) {
	cases := []TestCase{
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
			expected: map[string]int{
				"count": 2,
			},
		},
		// empty graph
		{
			input: []byte(`
			{
			}
			`),
			expected: map[string]int{
				"count": 0,
			},
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
			expected: map[string]int{
				"count": 1,
			},
		},
		// one component
		{
			input: []byte(`
			{
				"foo": []
			}
			`),
			expected: map[string]int{
				"count": 1,
			},
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
			expected: map[string]int{
				"count": 3,
			},
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
			expected: map[string]int{
				"count": 2,
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
			expected: map[string]int{
				"count": 3,
			},
		},
		// complex interconnected graph with four components
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
			expected: map[string]int{
				"count": 4,
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
			expected: map[string]int{
				"count": 1,
			},
		},
		// one cycle with the node connected to itself is still one component
		{
			input: []byte(`
			{
				"foo": ["foo"]
			}
			`),
			expected: map[string]int{
				"count": 1,
			},
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
			expected: map[string]int{
				"count": 2,
			},
		},
	}

	for _, testCase := range cases {
		MockReadFile := func(filePath string) ([]byte, error) {
			return testCase.input, nil
		}
		metricsItems := []string{MetricConnectedComponentsCount}
		result, err := metrics("mock.json", "", metricsItems, MockReadFile)
		if err != nil {
			t.Fail()
		}
		var actualOutput map[string]map[string]int
		json.Unmarshal(result, &actualOutput)
		assert.Equal(t, testCase.expected, actualOutput["components-count"])
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

	MockReadFile := func(filePath string) ([]byte, error) {
		return input, nil
	}

	metricsItems := []string{MetricDependenciesDirect, MetricDependenciesTransitive, MetricReverseDependenciesDirect, MetricReverseDependenciesTransitive}
	result, err := metrics("mock.json", "", metricsItems, MockReadFile)
	if err != nil {
		t.Fail()
	}
	var actualOutput map[string]map[string]int
	json.Unmarshal(result, &actualOutput)

	// check that all metrics are present in the report
	for _, metric := range metricsItems {
		_, exists := actualOutput[metric]
		assert.True(t, exists, "Expected metric '%s' to exist in the report", metric)
	}
}
