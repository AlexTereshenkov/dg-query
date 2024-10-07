#!/bin/bash

# Define the path to the genhtml binary inside the extracted lcov package
GENHTML_BINARY=$(bazel info output_base)/external/lcov/bin/genhtml

# Run Bazel coverage
bazel coverage //... --combined_report=lcov

# Generate the HTML report using genhtml from the downloaded package
$GENHTML_BINARY -o coverage-html bazel-out/_coverage/_coverage_report.dat
