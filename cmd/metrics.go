/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"log"
	"slices"
	"strings"
)

type GenericMapStringToAny map[string]interface{}

const (
	MetricDependenciesDirect            = "deps-direct"
	MetricDependenciesTransitive        = "deps-transitive"
	MetricReverseDependenciesDirect     = "rdeps-direct"
	MetricReverseDependenciesTransitive = "rdeps-transitive"
)

var allowedMetrics = []string{
	MetricDependenciesDirect,
	MetricDependenciesTransitive,
	MetricReverseDependenciesDirect,
	MetricReverseDependenciesTransitive,
}

func isValidMetric(metric string) bool {
	for _, allowedMetric := range allowedMetrics {
		if metric == allowedMetric {
			return true
		}
	}
	return false
}

/*
Get direct dependencies given adjacency list (used both for dependencies and dependents)
as it's just a mapping of node -> [nodes].
*/
func getMetricDependenciesDirect(adjacencyList AdjacencyList) GenericMapStringToAny {
	depsCount := make(GenericMapStringToAny)
	for key, deps := range adjacencyList {
		depsCount[key] = len(deps)
	}
	return depsCount
}

/*
Get transitive dependencies given adjacency list (used both for dependencies and dependents)
as it's just a mapping of node -> [nodes].
*/
func getMetricDependenciesTransitive(adjacencyList AdjacencyList) GenericMapStringToAny {
	depsCount := make(GenericMapStringToAny)
	visited := make(map[string]int)
	for dep := range adjacencyList {
		if count, found := visited[dep]; found {
			depsCount[dep] = count
		} else {
			transitiveDeps := getDepsTransitive(adjacencyList, []string{dep}, 0)
			depsCount[dep] = len(transitiveDeps)
			visited[dep] = len(transitiveDeps)
		}
	}
	return depsCount
}

// to be used in non-unit tests
var Metrics = metrics

/*
Produce data for given metrics.
*/
func metrics(filePathDg string, filePathDgReverse string, metricsItems []string, readFile ReadFileFunc) []byte {
	var adjacencyList AdjacencyList
	report := make(map[string]map[string]interface{})
	// use dependencies adjacency list as is
	if slices.Contains(metricsItems, MetricDependenciesDirect) || slices.Contains(metricsItems, MetricDependenciesTransitive) {
		jsonData := readFile(filePathDg)
		adjacencyList = loadJsonFile(jsonData)
	}
	// use the reversed dependencies list if provided otherwise reverse the dependencies list first
	if slices.Contains(metricsItems, MetricReverseDependenciesDirect) || slices.Contains(metricsItems, MetricReverseDependenciesTransitive) {
		if filePathDgReverse != "" {
			jsonData := readFile(filePathDgReverse)
			adjacencyList = loadJsonFile(jsonData)
		} else {
			jsonData := readFile(filePathDg)
			adjacencyList = reverseAdjacencyLists(loadJsonFile(jsonData))
		}
	}

	for _, metric := range metricsItems {
		if !isValidMetric(metric) {
			log.Printf("invalid metric: %s. Allowed metrics are: %s\n", metric, strings.Join(allowedMetrics, ","))
			return []byte("")
		}
		switch metric {

		case MetricDependenciesDirect:
			report[metric] = getMetricDependenciesDirect(adjacencyList)

		case MetricDependenciesTransitive:
			report[metric] = getMetricDependenciesTransitive(adjacencyList)

		case MetricReverseDependenciesDirect:
			report[metric] = getMetricDependenciesDirect(adjacencyList)

		case MetricReverseDependenciesTransitive:
			report[metric] = getMetricDependenciesTransitive(adjacencyList)
		}
	}
	reportJson, _ := json.MarshalIndent(report, "", "  ")
	return reportJson
}
