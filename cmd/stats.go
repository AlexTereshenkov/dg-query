/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
List dependencies for given targets.
*/
func stats(cmd *cobra.Command, metrics []string) {
	filePath, _ := cmd.Flags().GetString("dg")
	jsonData := ReadFile(filePath)
	adjacencyList := loadJsonFile(jsonData)
	fmt.Fprintln(cmd.OutOrStdout(), adjacencyList)

}
