# Build

This document provides an overview of the build process. 

## Bazel

Bazel is used to build the project locally and in CI. However, the project can still
be built with native Go tooling. For local builds, run `make build-go` or `make build-bazel`
to build with Go and Bazel, respectively. The project originally used `WORKSPACE` based build, 
but then was migrated to Bzlmod.

For CI, a few scripts to be run in Buildkite exist.

* `.buildkite/build-go.sh`: build on Buildkite agent using native Go tooling
* `.buildkite/build-bazel.sh`: build on Buildkite agent using Bazel
* `.buildkite/build-bazel-docker.sh`: build in a Docker container (using `--spawn_strategy=docker`, see the [docs](https://bazel.build/reference/command-line-reference#build-flag--spawn_strategy))
* `.buildkite/build-bazel-remote.sh`: build remotely on BuildBuddy (on their agent)

The `defs` directory contains example of macros and loading members of `.bzl` files.

The `platforms` directory contains examples on how to:
* build CGo code
* run tests on Windows 
* build in a custom Docker container

Build specific documentation can be found at the `.buildkite` directory's Shell scripts 
and files with build metadata such as `MODULE.bazel` and `BUILD.bazel` files.

## Release

A GitHub release is produced by running the following commands

```shell
git tag v0.X.0
git push --tags
```

which triggers the `goreleaser` action in `.github/workflows/release.yml`.