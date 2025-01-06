//go:build windows

/*
Copyright © 2024 Alexey Tereshenkov
*/
package cmd

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
This test is only to be used to confirm that Bazel platforms
in the project are configured correctly; it is only supposed to
be run on a Windows device.
*/
func TestWindows(t *testing.T) {
	assert.Equal(t, "windows", runtime.GOOS)
}
