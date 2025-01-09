#!/usr/bin/env bash
# adopted from https://github.com/buildbuddy-io/buildbuddy/blob/master/workspace_status.sh

# This script will be run bazel when building process starts to
# generate key-value information that represents the status of the
# workspace. The output should be like
#
# KEY1 VALUE1
# KEY2 VALUE2
#
# If the script exits with non-zero code, it's considered as a failure
# and the output will be discarded.

# these are sent to BuildBuddy as build metadata
# it should be `GIT_BRANCH` as per
# https://github.com/buildbuddy-io/buildbuddy/blob/master/workspace_status.sh
# when called from the workspace status Shell script, and it should be
# BRANCH_NAME as per https://www.buildbuddy.io/docs/guide-metadata#build-metadata-1
# when called from the command line.
if [[ -n "$BUILDKITE_BRANCH" ]]; then
    git_branch="${BUILDKITE_BRANCH}"
else
    git_branch="$(git rev-parse --abbrev-ref HEAD)"
fi
echo "GIT_BRANCH $git_branch"

echo "COMMIT_SHA $(git rev-parse HEAD)"

echo "REPO_URL https://github.com/AlexTereshenkov/dg-query.git"

# The "STABLE_" suffix causes these to be part of the "stable" workspace
# status, which may trigger rebuilds of certain targets if these values change
# and you're building with the "--stamp" flag.

# this is embedded into the Go release binary
echo "STABLE_GIT_COMMIT $(git rev-parse HEAD)"
