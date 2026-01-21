#!/bin/bash
# This script is supposed to run after the `bazel coverage //... --combined_report=lcov` command
# as it expects `bazel-out/_coverage/_coverage_report.dat` to exist on disk and it is not recommended
# to run Bazel commands from within Bazel.
# Pure Bazel approach would be:
# bazel coverage --instrument_test_targets --combined_report=lcov //...
# bazel run @lcov//:genhtml -- --source-directory $(pwd) --branch-coverage --output $(pwd)/genhtml \
#   $(bazel info output_path)/_coverage/_coverage_report.dat
echo "Running codecov.sh Shell script to generate the HTML code coverage report"

# Define the path to the genhtml binary inside the extracted lcov package
GENHTML_BINARY=$1
# Generate the HTML report using genhtml from the downloaded package
$GENHTML_BINARY -o $BUILD_WORKSPACE_DIRECTORY/coverage-html $BUILD_WORKSPACE_DIRECTORY/bazel-out/_coverage/_coverage_report.dat
