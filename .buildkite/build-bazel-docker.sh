#! /bin/bash
set -xueo pipefail

echo '--- Build'
bazel --version
build-support/create-bazelrc-ci.sh

# `--load` ensures that the image built is loaded into the local Docker environment, 
# making it immediately available to run or reference by tag; otherwise, the image 
# will only exist in Docker's build cache.
docker build -t go-build-env -f platforms/Dockerfile --load .
docker image ls

bazel build --config="docker" //... 

echo '--- Test'
bazel test --config="docker" //...  | sed 's/^=== //' | sed 's/^--- //'

echo '--- Running a compiled CGo binary'
bazel run --config="docker" //platforms/compiled:cgo-compiled

echo '--- Finish'
