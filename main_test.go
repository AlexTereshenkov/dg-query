/*
Copyright © 2024 Alexey Tereshenkov
*/
package main

import (
	"bytes"
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
				"bla.py"
			],
			"spam.py": [
				"eggs.py",
				"ham.py"
			]
		
		}		
		`)
	}
	cmd.RootCmd.SetArgs([]string{"dependencies", "--dg=dg.json", "foo.py"})
	cmd.RootCmd.Execute()

	expected := "bar.py\nbla.py\n"
	actualOutput := buf.String()
	assert.Equal(t, expected, actualOutput, "Failing assertion")
}
