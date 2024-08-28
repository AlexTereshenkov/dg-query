/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AlexTereshenkov/dg-query/cmd"
	"github.com/stretchr/testify/assert"
)

func TestDependencies(t *testing.T) {

	var buf bytes.Buffer
	// redirection is not required for any subcommands, but this is how it's done for the reference:
	// for _, c := range cmd.RootCmd.Commands() {
	// 	c.SetOut(&buf)
	// }
	cmd.RootCmd.SetOut(&buf)
	cmd.RootCmd.SetErr(&buf)

	// mocking function that reads a file from disk
	cmd.ReadFile = func(filePath string) []byte {
		return []byte(`
		{
			"foo.py": [
				"bar.py",
				"baz.py"
			],
			"spam.py": [
				"eggs.py",
				"ham.py"
			]
		
		}		
		`)
	}
	var expected []string

	cmd.RootCmd.SetArgs([]string{"dependencies", "--dg=dg.json", "foo.py"})
	cmd.RootCmd.Execute()
	expected = []string{"bar.py", "baz.py"}
	actualOutput := strings.Split(buf.String(), "\n")[:2]
	assert.ElementsMatch(t, expected, actualOutput, "Failing assertion")
	buf.Reset()

	// asking for non-existing node
	cmd.RootCmd.SetArgs([]string{"dependencies", "--dg=dg.json", "non-existing.py"})
	cmd.RootCmd.Execute()
	expected = []string{"\n"}
	actualOutput = strings.Split(buf.String(), "")
	assert.ElementsMatch(t, expected, actualOutput, "Failing assertion")
	buf.Reset()

	// asking for multiple nodes
	cmd.RootCmd.SetArgs([]string{"dependencies", "--dg=dg.json", "foo.py", "spam.py"})
	cmd.RootCmd.Execute()
	expected = []string{"bar.py", "baz.py", "eggs.py", "ham.py"}
	actualOutput = strings.Split(buf.String(), "\n")[:4]
	assert.ElementsMatch(t, expected, actualOutput, "Failing assertion")
	buf.Reset()
}
