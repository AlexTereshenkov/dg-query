/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands.
// This needs to be exported to be available to `main.go`.
var RootCmd = &cobra.Command{
	Use:   "dg-query",
	Short: "A command-line utility program to query dependency graph of a codebase.",
	Long: `A command-line utility program to query dependency graph of a codebase
which operates on the adjacency list data stored as a JSON file.`,
}

// getting dependencies of a particular target in the dependency graph
var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Get dependencies of a target",
	Long:  `Get dependencies of a target`,
	// requiring at least one target address (to get their dependencies)
	Args: cobra.MinimumNArgs(1),
	Run:  dependencies,
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get dependency graph related statistics",
	Long:  `Get dependency graph related statistics`,
	Run:   stats,
}

// JSON file with the dependency graph represented as an adjacency list
var dg string

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
	RootCmd.AddCommand(statsCmd)
	RootCmd.AddCommand(dependenciesCmd)

	//make dg flag global for all commands as all of them will need dg data
	// TODO: make flag obligatory
	RootCmd.PersistentFlags().StringVar(&dg, "dg", "", "JSON file with the dependency graph represented as an adjacency list")

	// TODO: what is more common to pass as --metrics="deps_count,rdeps_count" 
	// or as --metric="deps_count" --metric="rdeps_count"

	dependenciesCmd.Flags().BoolP("transitive", "", false, "Get transitive dependencies")
	dependenciesCmd.Flags().BoolP("reflexive", "", false, "Include input targets in the output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
