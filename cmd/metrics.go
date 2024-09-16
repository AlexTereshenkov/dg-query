/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

type MetricDependenciesDirectType map[string]interface{}
type MetricDependenciesTransitiveType map[string]interface{}

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

func getMetricDependenciesDirect(adjacencyList AdjacencyList) MetricDependenciesDirectType {
	depsCount := make(MetricDependenciesDirectType)
	for key, deps := range adjacencyList {
		depsCount[key] = len(deps)
	}
	return depsCount
}

func getMetricDependenciesTransitive(adjacencyList AdjacencyList) MetricDependenciesTransitiveType {
	depsCount := make(MetricDependenciesTransitiveType)
	visited := make(map[string]int)
	for dep := range adjacencyList {
		if count, found := visited[dep]; found {
			depsCount[dep] = count
		} else {
			transitiveDeps := getDepsTransitive(adjacencyList, []string{dep})
			depsCount[dep] = len(transitiveDeps)
			visited[dep] = len(transitiveDeps)
		}
	}
	return depsCount
}

/*
Produce data for given metrics.
*/
func metrics(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("dg")
	jsonData := ReadFile(filePath)
	adjacencyList := loadJsonFile(jsonData)
	metricsItems, _ := cmd.Flags().GetStringSlice("metric")

	report := make(map[string]map[string]interface{})

	for _, metric := range metricsItems {
		if !isValidMetric(metric) {
			log.Printf("invalid metric: %s. Allowed metrics are: %s\n", metric, strings.Join(allowedMetrics, ","))
			return
		}
		switch metric {
		case MetricDependenciesDirect:
			report[metric] = getMetricDependenciesDirect(adjacencyList)

		case MetricDependenciesTransitive:
			report[metric] = getMetricDependenciesTransitive(adjacencyList)
		}

	}
	reportJson, _ := json.MarshalIndent(report, "", "  ")
	cmd.OutOrStdout().Write(reportJson)
	cmd.OutOrStdout().Write([]byte("\n"))

}
