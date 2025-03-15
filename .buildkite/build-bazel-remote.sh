#! /bin/bash -xue

echo '--- Build'
bazel --version
build-support/create-bazelrc-ci.sh
echo "build --remote_executor=grpcs://remote.buildbuddy.io" >> .bazelrc.ci
bazel build //... 

echo '--- Test'
bazel test //... | sed 's/^=== //' | sed 's/^--- //'

echo '--- Running a compiled CGo binary'
bazel run //platforms/compiled:cgo-compiled

echo '--- Finish'
