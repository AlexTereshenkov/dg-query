#!/usr/bin/env bash

echo "common --show_timestamps" > .bazelrc.ci
echo "common --config=bb" >> .bazelrc.ci
echo "build --build_metadata=ROLE=CI" >> .bazelrc.ci

# this secret is declared in Buildkite secrets
if command -v buildkite-agent &> /dev/null
then
    # runs in Buildkite
    BUILD_BUDDY_API_KEY=$(buildkite-agent secret get BUILD_BUDDY_API_KEY)
else
    # runs locally and the api key should be available via the environment variable
    BUILD_BUDDY_API_KEY=${BUILD_BUDDY_API_KEY}
fi
echo "build:bb --remote_header=x-buildbuddy-api-key=$BUILD_BUDDY_API_KEY" >> .bazelrc.ci
