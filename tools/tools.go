/*
This file lists tools as direct dependencies even though
they are imported in Go source code; this is to ensure
`bazel mod tidy` and `go mod tidy` are not in conflict.
This is necessary since `go.mod` tracks direct dependencies
(so we add dummy `import package-name` in a `.go` file).
*/

//go:build tools
package tools

import (
    _ "golang.org/x/tools/cmd/goimports"
)
