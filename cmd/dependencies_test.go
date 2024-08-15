/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/cobra"

)

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return output, err
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestDependencies(t *testing.T) {
	// var buf bytes.Buffer
	// writer := &buf
	// stdout := os.Stdout
	// os.Stdout = writer
	// defer func() { os.Stdout = stdout }()

	// Dependencies

	// cmd := NewrootCmd()
	// consoleOutput, err := ExecuteCommand(cmd, "small")
	// assert.Equal(consoleOutput, "")

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"A", "a1"})
	rootCmd.Execute()

	expected := "This-is-command-a1"

	assert.Equal(t, actual.String(), expected, "actual is not expected")
}
