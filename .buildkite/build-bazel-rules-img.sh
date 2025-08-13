#! /bin/bash
set -xueo pipefail

echo '--- Build'
bazel --version
build-support/create-bazelrc-ci.sh

echo '--- Authenticate with Docker registry'
bazel run @tweag-credential-helper//installer

echo '--- Start Docker registry'
docker run -d -p 5000:5000 --name registry registry:2
curl http://localhost:5000/v2/

echo '--- Build Docker image'
bazel build //:app_image

echo '--- Push Docker image to local registry'
bazel run //:push

echo '--- Check the list of images'
curl http://localhost:5000/v2/repository/my-project/tags/list

echo '--- Run the app'
docker run --rm localhost:5000/repository/my-project:latest bash -c "/app/bin/dg-query"

echo '--- Finish'
