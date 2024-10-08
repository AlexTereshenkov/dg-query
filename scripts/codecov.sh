#!/bin/bash

# Define the path to the genhtml binary inside the extracted lcov package
GENHTML_BINARY=$1 #(bazel info output_base)/external/lcov/bin/genhtml

# Generate the HTML report using genhtml from the downloaded package
$GENHTML_BINARY -o coverage-html $BUILD_WORKSPACE_DIRECTORY/bazel-out/_coverage/_coverage_report.dat
