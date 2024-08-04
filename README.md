# Environment setup

## Local development

Documenting what was done during the project bootstrap. These steps are not required to run.

```shell
go env -w GOPATH=$HOME/go
vi ~/.zshrc
export PATH="$PATH:/home/$USER/.local/bin:$GOPATH/bin"

go install github.com/spf13/cobra-cli@latest
go mod init dg-query
go mod tidy # to update `go.mod` file, if necessary
cobra-cli init
go build
go install # put `dg-query` binary into the `$GOPATH/bin` directory
dg-query # should run the app
```

## Bazel

These are the commands one would run to build and run the app.

```shell
# download `buildifier` into the `$GOPATH/bin` directory, if not installed, and reformat
go install github.com/bazelbuild/buildtools/buildifier@latest
buildifier -r .

# update 3rd-party dependencies
gazelle update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies -prune

# update local BUILD files
gazelle

# run the app
bazel build //... && bazel-bin/dg-query_/dg-query

OR

bazel run //:dg-query
```

## CI

```shell
gofmt -w .
go mod tidy
buildifier -r .
```