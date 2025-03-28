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
	MetricConnectedComponentsCount      = "components-count"
)

var allowedMetrics = []string{
	MetricDependenciesDirect,
	MetricDependenciesTransitive,
	MetricReverseDependenciesDirect,
	MetricReverseDependenciesTransitive,
	MetricConnectedComponentsCount,
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
	// set count to 0 for nodes that were not present as keys in the adjacency list
	// (i.e. they are only dependencies of some nodes and do not depend on anything)
	for _, deps := range adjacencyList {
		for _, dep := range deps {
			if _, isKey := adjacencyList[dep]; !isKey {
				depsCount[dep] = 0
			}
		}
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
	// set count to 0 for nodes that were not present as keys in the adjacency list
	// (i.e. they are only dependencies of some nodes and do not depend on anything)
	for _, deps := range adjacencyList {
		for _, dep := range deps {
			if _, isKey := adjacencyList[dep]; !isKey {
				depsCount[dep] = 0
			}
		}
	}
	return depsCount
}

// getConnectedComponentsCount gets count of connected components in a graph
func getConnectedComponentsCount(adjacencyList AdjacencyList) GenericMapStringToAny {
	connectedComponentsCount := make(GenericMapStringToAny)
	connectedComponentsCount["count"] = len(getConnectedComponents(adjacencyList))
	return connectedComponentsCount
}

// to be used in non-unit tests
var Metrics = metrics

/*
Produce data for given metrics.
*/
func metrics(filePathDg string, filePathDgReverse string, metricsItems []string, readFile ReadFileFunc) ([]byte, error) {
	var adjacencyList AdjacencyList
	var adjacencyListReverse AdjacencyList

	report := make(map[string]map[string]interface{})
	// use dependencies adjacency list as is
	if slices.Contains(metricsItems, MetricDependenciesDirect) ||
		slices.Contains(metricsItems, MetricDependenciesTransitive) ||
		slices.Contains(metricsItems, MetricConnectedComponentsCount) {
		jsonData, err := readFile(filePathDg)
		if err != nil {
			return nil, err
		}
		adjacencyList, err = loadJsonFile(jsonData)
		if err != nil {
			return nil, err
		}
	}
	// use the reversed dependencies list if provided otherwise reverse the dependencies list first
	if slices.Contains(metricsItems, MetricReverseDependenciesDirect) || slices.Contains(metricsItems, MetricReverseDependenciesTransitive) {
		if filePathDgReverse != "" {
			jsonData, err := readFile(filePathDgReverse)
			if err != nil {
				return nil, err
			}
			adjacencyListReverse, err = loadJsonFile(jsonData)
			if err != nil {
				return nil, err
			}
		} else {
			jsonData, err := readFile(filePathDg)
			if err != nil {
				return nil, err
			}
			adjacencyListTemp, err := loadJsonFile(jsonData)
			if err != nil {
				return nil, err
			}
			adjacencyListReverse = reverseAdjacencyLists(adjacencyListTemp)
			// extending the map with the nodes that had no dependencies (e.g. {"foo": []} as "foo" won't be in reverse adjacency list
			// if no one depends on "foo")
			for key, deps := range adjacencyListTemp {
				if len(deps) == 0 {
					adjacencyListReverse[key] = []string{}
				}
			}
		}
	}

	for _, metric := range metricsItems {
		if !isValidMetric(metric) {
			log.Printf("invalid metric: %s. Allowed metrics are: %s\n", metric, strings.Join(allowedMetrics, ","))
			return []byte(""), nil
		}
		switch metric {

		case MetricDependenciesDirect:
			report[metric] = getMetricDependenciesDirect(adjacencyList)

		case MetricDependenciesTransitive:
			report[metric] = getMetricDependenciesTransitive(adjacencyList)

		case MetricReverseDependenciesDirect:
			report[metric] = getMetricDependenciesDirect(adjacencyListReverse)

		case MetricReverseDependenciesTransitive:
			report[metric] = getMetricDependenciesTransitive(adjacencyListReverse)

		case MetricConnectedComponentsCount:
			report[metric] = getConnectedComponentsCount(adjacencyList)

		}
	}
	reportJson, _ := json.MarshalIndent(report, "", "  ")
	return reportJson, nil
}
