#!/bin/bash
# This script is supposed to run after the `bazel coverage //... --combined_report=lcov` command
# as it expects `bazel-out/_coverage/_coverage_report.dat` to exist on disk and it is not recommended
# to run Bazel commands from within Bazel
echo "Running codecov.sh Shell script to generate the HTML code coverage report"

# Define the path to the genhtml binary inside the extracted lcov package
GENHTML_BINARY=$1

# Generate the HTML report using genhtml from the downloaded package
RUNFILES_ROOT=$PWD
$GENHTML_BINARY -o $BUILD_WORKSPACE_DIRECTORY/coverage-html $BUILD_WORKSPACE_DIRECTORY/bazel-out/_coverage/_coverage_report.dat
