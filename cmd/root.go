/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// This variable is going to be set by Bazel during stamping a Go binary;
// https://github.com/bazel-contrib/rules_go/blob/master/docs/go/core/defines_and_stamping.md
var Version = "(redacted)"

// RootCmd represents the base command when called without any subcommands.
// This needs to be exported to be available to `main.go`.
var RootCmd = &cobra.Command{
	Use:   "dg-query",
	Short: "A command-line utility program to query dependency graph of a codebase.",
	Long: `A command-line utility program to query dependency graph of a codebase
which operates on the adjacency list data stored as a JSON file. 

Git revision: ` + Version,
}

// getting dependencies of given targets in the dependency graph
var dependenciesCmd = &cobra.Command{
	Use:     "dependencies",
	Aliases: []string{"deps"},
	Short:   "Get dependencies of given targets",
	Long:    `Get dependencies of given targets`,
	// requiring at least one target address (to get their dependencies)
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")
		transitive, _ := cmd.Flags().GetBool("transitive")
		reflexive, _ := cmd.Flags().GetBool("reflexive")
		depth, _ := cmd.Flags().GetInt("depth")

		result, err := dependencies(filePath, targets, transitive, reflexive, depth, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output := strings.Join(result, "\n")
		cmd.OutOrStdout().Write([]byte(output))
		cmd.OutOrStdout().Write([]byte("\n"))
	},
}

// getting all roots in the dependency graph
var rootsCmd = &cobra.Command{
	Use:   "roots",
	Short: "Get nodes that no other node depends on",
	Long:  `Get nodes that no other node depends on`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")

		result, err := roots(filePath, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output := strings.Join(result, "\n")
		cmd.OutOrStdout().Write([]byte(output))
		cmd.OutOrStdout().Write([]byte("\n"))
	},
}

// getting all leaves in the dependency graph
var leavesCmd = &cobra.Command{
	Use:   "leaves",
	Short: "Get nodes that have no dependencies",
	Long:  `Get nodes that have no dependencies`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")

		result, err := leaves(filePath, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output := strings.Join(result, "\n")
		cmd.OutOrStdout().Write([]byte(output))
		cmd.OutOrStdout().Write([]byte("\n"))
	},
}

// getting dependents of given targets in the dependency graph
var dependentsCmd = &cobra.Command{
	Use:     "dependents",
	Aliases: []string{"rdeps"},
	Short:   "Get dependents of given targets",
	Long:    `Get dependents of given targets`,
	// requiring at least one target address (to get their dependents)
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, targets []string) {
		filePathDg, _ := cmd.Flags().GetString("dg")
		filePathDgReverse, _ := cmd.Flags().GetString("rdg")
		transitive, _ := cmd.Flags().GetBool("transitive")
		reflexive, _ := cmd.Flags().GetBool("reflexive")
		depth, _ := cmd.Flags().GetInt("depth")

		result, err := dependents(filePathDg, filePathDgReverse,
			targets, transitive, reflexive, depth, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output := strings.Join(result, "\n")
		cmd.OutOrStdout().Write([]byte(output))
		cmd.OutOrStdout().Write([]byte("\n"))
	},
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Get dependency graph related metrics",
	Long:  `Get dependency graph related metrics`,
	Run: func(cmd *cobra.Command, args []string) {
		filePathDg, _ := cmd.Flags().GetString("dg")
		filePathDgReverse, _ := cmd.Flags().GetString("rdg")
		metricsItems, _ := cmd.Flags().GetStringSlice("metric")
		result, err := metrics(filePathDg, filePathDgReverse, metricsItems, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cmd.OutOrStdout().Write(result)
		cmd.OutOrStdout().Write([]byte("\n"))

	},
}

var pathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Get paths between targets",
	Long:  `Get paths between targets`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")
		fromTarget, _ := cmd.Flags().GetString("from")
		toTarget, _ := cmd.Flags().GetString("to")
		maxPaths, _ := cmd.Flags().GetInt("n")
		result, err := paths(filePath, fromTarget, toTarget, maxPaths, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resultJson, _ := json.MarshalIndent(result, "", "  ")
		cmd.OutOrStdout().Write(resultJson)
		cmd.OutOrStdout().Write([]byte("\n"))

	},
}

var cyclesCmd = &cobra.Command{
	Use:   "cycles",
	Short: "Find cycles in the dependency graph",
	Long:  `Find cycles in the dependency graph`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")
		result, err := cycles(filePath, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resultJson, _ := json.MarshalIndent(result, "", "  ")
		cmd.OutOrStdout().Write(resultJson)
		cmd.OutOrStdout().Write([]byte("\n"))

	},
}

var subgraphCmd = &cobra.Command{
	Use:   "subgraph",
	Short: "Extract a subgraph out of the dependency graph",
	Long:  `Extract a subgraph out of the dependency graph`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")
		rootNode, _ := cmd.Flags().GetString("root")
		result, err := extractSubgraph(filePath, rootNode, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resultJson, _ := json.MarshalIndent(result, "", "  ")
		cmd.OutOrStdout().Write(resultJson)
		cmd.OutOrStdout().Write([]byte("\n"))

	},
}

var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Get a list of connected components in the dependency graph",
	Long:  `Get a list of connected components in the dependency graph`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")
		result, err := listConnectedComponents(filePath, DefaultReadFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resultJson, _ := json.MarshalIndent(result, "", "  ")
		cmd.OutOrStdout().Write(resultJson)
		cmd.OutOrStdout().Write([]byte("\n"))

	},
}

var simplifyCmd = &cobra.Command{
	Use:   "simplify",
	Short: "Simplify the dependency graph by applying a requested technique",
	Long:  `Simplify the dependency graph by applying a requested technique`,
	Run: func(cmd *cobra.Command, targets []string) {
		filePath, _ := cmd.Flags().GetString("dg")
		technique, _ := cmd.Flags().GetString("technique")
		result, err := simplifyAdjacencyList(filePath, DefaultReadFile, technique)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		resultJson, _ := json.MarshalIndent(result, "", "  ")
		cmd.OutOrStdout().Write(resultJson)
		cmd.OutOrStdout().Write([]byte("\n"))

	},
}

// JSON file with the dependency graph represented as an adjacency list
var dg string
var rdg string

// metrics to be generated by the `metrics` command
var metricsFlags []string

// paths command from and to targets
var fromTarget string
var toTarget string

// subgraph command root node
var rootNode string

// simplify command technique to apply
var simplifyTechnique string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.AddCommand(metricsCmd)
	RootCmd.AddCommand(pathsCmd)
	RootCmd.AddCommand(cyclesCmd)
	RootCmd.AddCommand(subgraphCmd)
	RootCmd.AddCommand(dependenciesCmd)
	RootCmd.AddCommand(dependentsCmd)
	RootCmd.AddCommand(componentsCmd)
	RootCmd.AddCommand(rootsCmd)
	RootCmd.AddCommand(leavesCmd)
	RootCmd.AddCommand(simplifyCmd)

	//make dg flag global for all commands as all of them will need dg data
	RootCmd.PersistentFlags().StringVar(&dg, "dg", "", "JSON file with the dependency graph represented as an adjacency list")

	pathsCmd.Flags().StringVar(&dg, "dg", "", "JSON file with the dependency graph represented as an adjacency list")
	pathsCmd.Flags().StringVar(&fromTarget, "from", "", "Find path from this target")
	pathsCmd.Flags().StringVar(&toTarget, "to", "", "Find path to this target")
	pathsCmd.Flags().Int("n", 0, "Only return first n paths between targets")

	metricsCmd.Flags().StringVar(&rdg, "rdg", "", "JSON file with the dependency graph represented as an adjacency list")
	metricsCmd.Flags().StringSliceVar(&metricsFlags, "metric", []string{}, "Metrics to report")

	dependenciesCmd.Flags().BoolP("transitive", "", false, "Get transitive dependencies")
	dependenciesCmd.Flags().BoolP("reflexive", "", false, "Include input targets in the output")
	dependenciesCmd.Flags().Int("depth", 0, "Depth of search for transitive dependencies")

	dependentsCmd.Flags().StringVar(&rdg, "rdg", "", "JSON file with the dependency graph represented as an adjacency list")
	dependentsCmd.Flags().BoolP("transitive", "", false, "Get transitive dependents")
	dependentsCmd.Flags().BoolP("reflexive", "", false, "Include input targets in the output")
	dependentsCmd.Flags().Int("depth", 0, "Depth of search for transitive dependents")

	subgraphCmd.Flags().StringVar(&rootNode, "root", "", "Root node for the subgraph to extract")

	simplifyCmd.Flags().StringVar(&simplifyTechnique, "technique", "", "Technique to simplify the dependency graph")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
