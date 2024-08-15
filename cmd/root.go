/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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
	Run:  Dependencies,
}

// JSON file with the dependency graph represented as an adjacency list
var dg string
var target string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.AddCommand(dependenciesCmd)
	dependenciesCmd.Flags().StringVar(&dg, "dg", "", "JSON file with the dependency graph represented as an adjacency list")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
