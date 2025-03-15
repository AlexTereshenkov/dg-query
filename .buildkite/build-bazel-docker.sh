#! /bin/bash -xue

echo '--- Build'
bazel --version
build-support/create-bazelrc-ci.sh
# `--load` ensures that the image built is loaded into the local Docker environment, 
# making it immediately available to run or reference by tag; otherwise, the image 
# will only exist in Docker's build cache.
docker build -t go-build-env -f platforms/Dockerfile --load .
docker image ls
bazel build //... --spawn_strategy=docker --experimental_enable_docker_sandbox --platforms="//platforms:docker"

echo '--- Test'
bazel test //...  --spawn_strategy=docker --experimental_enable_docker_sandbox --platforms="//platforms:docker" | sed 's/^=== //' | sed 's/^--- //'

echo '--- Running a compiled CGo binary'
bazel run //platforms/compiled:cgo-compiled

echo '--- Finish'
